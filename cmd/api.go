package cmd

import (
	"bufio"
	"crypto/sha256"
	"device-go/crypto"
	"device-go/shared/config"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"

	log "github.com/ndmsystems/golog"
)

// CMD is a command to execute
func Build(b []byte) (CMD, error) {
	c := CMD{}
	err := json.Unmarshal(b, &c)
	return c, err
}

// convert cmd to readable string with limit by body by 50
func (c CMD) Readable() string {
	body := c.Body
	if len(body) > 50 {
		body = body[:50]
	}
	return "uuid: " + c.UUID + " ts: " + string(c.Ts) + " sender: " + c.Sender + " sender_nacl: " + c.SenderNacl + " receiver: " + string(c.Receiver) + " hash: " + c.Hash + " sign: " + c.Sign + " body: " + body
}

// Valid checks if all fields are filled
func (c CMD) Valid() bool {
	log.Debug(c.UUID)

	if len(c.UUID) == 0 {
		return false
	}
	if c.Ts == 0 {
		return false
	}
	if len(c.Sender) == 0 {
		return false
	}
	if len(c.SenderNacl) == 0 {
		return false
	}
	if len(c.Receiver) == 0 {
		return false
	}
	if len(c.Hash) == 0 {
		return false
	}
	if len(c.Sign) == 0 {
		return false
	}
	if len(c.Body) == 0 {
		return false
	}
	return true
}

// GetHash calculates hash sum of all fields except Sign and Hash
func (c CMD) GetHash() string {
	log.Debug(c.UUID)
	h := sha256.New()
	bt := []byte(c.UUID + string(c.Ts) + c.Sender + c.SenderNacl + string(c.Receiver) + c.Body)
	h.Write(bt)
	return string(h.Sum(nil))
}

// check signature of command result of verification with public key of sender
func (c CMD) Verify() (string, bool) {
	log.Debug(c.UUID)
	if !config.IsProd() || c.Sign != "testing" { // for testing purposes only "testing" signature is allowed
		return c.Sender, true
	}
	return crypto.Keys.Verify(c.Sign)
}

// Execute executes command and returns result and error if any occurs
func (c CMD) Execute() (string, error) {

	log.Debug(c.Readable())

	body, err := crypto.Keys.Decrypt(c.Body, c.SenderNacl)
	if err != nil {
		return "", err
	}

	cmdArr := strings.Split(body, " ")
	var args []string
	if len(cmdArr) > 1 {
		args = cmdArr[1:]
	}
	log.Debug(cmdArr[0], args)
	cmd := exec.Command(cmdArr[0], args...)
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}
	log.Debug(cmd.String(), args)

	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		log.Error(err)
		return "", err
	}

	var errOut string
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		errOut += scanner.Text() + "\n"
	}
	if len(errOut) > 0 {
		log.Error(errOut)
		return "", errors.New(errOut)
	}

	var out string
	scanner = bufio.NewScanner(stdout)
	for scanner.Scan() {
		out += scanner.Text() + "\n"
	}
	log.Debug(out)
	return out, nil
}
