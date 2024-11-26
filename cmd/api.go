package cmd

import (
	"bytes"
	"crypto/sha256"
	"device-go/crypto"
	"device-go/shared/config"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"os/exec"
	"strconv"
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

// getHash calculates hash sum of all fields except Sign and Hash
func (c CMD) getHash() []byte {
	log.Debug(c.UUID)
	h := sha256.New()
	bt := []byte(c.UUID + strconv.FormatInt(c.Ts, 10) + c.Sender + c.SenderNacl + string(c.Receiver) + c.Body)
	h.Write(bt)
	return h.Sum(nil)
}

// VerifySignature of command result with public key of sender
func (c CMD) VerifySignature() bool {
	log.Debug(c.UUID)
	if !config.IsProd() {
		// for testing purposes "testing" signature is allowed as well as valid signature
		return c.Sign == "testing" || crypto.Keys.VerifySignature(c.Sender, c.getHash(), c.Sign)
	}
	return crypto.Keys.VerifySignature(c.Sender, c.getHash(), c.Sign)
}

func (c CMD) Execute() (string, error) {
	log.Debug(c.Readable())

	cmd, err := crypto.Keys.Decrypt(c.Body, c.SenderNacl)
	if err != nil {
		return "", err
	}
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return "", errors.New("no command provided")
	}

	// Осуществляет выполнение команды с сохранением форматирования вывода
	command := exec.Command(parts[0], parts[1:]...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &out
	command.Stderr = &stderr
	err = command.Run()
	if err != nil {
		return fmt.Sprintf("%s\n%s", out.String(), stderr.String()), err
	}
	return out.String(), nil
}
