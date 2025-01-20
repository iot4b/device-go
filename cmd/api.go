package cmd

import (
	"crypto/sha256"
	"device-go/crypto"
	"device-go/shared/config"
	"encoding/json"
	"fmt"
	"github.com/google/shlex"
	log "github.com/ndmsystems/golog"
	"github.com/pkg/errors"
	"os/exec"
	"strconv"
	"strings"
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
	return fmt.Sprintf("uuid: %s ts: %v sender: %s sender_nacl: %s receiver: %s hash: %s sign: %s body: %s",
		c.UUID, c.Ts, c.Sender, c.SenderNacl, c.Receiver, c.Hash, c.Sign, body)
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

	return run(cmd)
}

func run(cmd string) (string, error) {
	log.Debug("CMD:", cmd)

	parts, err := shlex.Split(cmd)
	if err != nil {
		return "", fmt.Errorf("error parsing command: %v", err)
	}
	if len(parts) == 0 {
		return "", errors.New("no command provided")
	}

	if parts[0] == "ndms" {
		// run keenetic command
		kcmd := fmt.Sprintf("ndmq -p \"%s\" -x", strings.Join(parts[1:], " "))
		log.Debug("Run Keenetic CMD:", kcmd)
		return run(kcmd)
	}

	// Осуществляет выполнение команды с сохранением форматирования вывода
	out, err := exec.Command(parts[0], parts[1:]...).CombinedOutput()
	return string(out), err
}
