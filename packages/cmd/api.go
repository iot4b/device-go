package cmd

import (
	"crypto/sha256"
	"device-go/packages/config"
	"device-go/packages/crypto"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/google/shlex"
	log "github.com/ndmsystems/golog"
	"github.com/pkg/errors"
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
	return fmt.Sprintf("uuid: %s ts: %v sender: %s receiver: %s hash: %s sign: %s body: %s",
		c.UUID, c.Ts, c.Sender, c.Receiver, c.Hash, c.Sign, body)
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
	bt := []byte(c.UUID + strconv.FormatInt(c.Ts, 10) + c.Sender + string(c.Receiver) + c.Body)
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

	// decrypt plain cmd
	cmd, err := crypto.Keys.DecryptChaCha20Poly1305(c.Body, c.Sender)
	if err != nil {
		return "", err
	}

	return c.run(cmd)
}

func (c CMD) run(cmd string) (string, error) {
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
		return c.run(kcmd)
	}

	// Осуществляет выполнение команды с сохранением форматирования вывода
	out, err := exec.Command(parts[0], parts[1:]...).CombinedOutput()
	log.Info("CMD:", cmd)

	// encrypt the response
	return crypto.Keys.EncryptChaCha20Poly1305(out, c.Sender)
}
