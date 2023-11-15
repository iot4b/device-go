package handlers

import (
	"bufio"
	"device-go/crypto"
	"device-go/dsm"
	"device-go/shared/config"
	"device-go/storage"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
)

type info struct {
	Address string `json:"address"`
	Version string `json:"version"`
	Elector string `json:"elector"`
	Node    string `json:"node"`
	Type    string `json:"type"`
}

// info для коалы
func GetInfo(_ *coalago.CoAPMessage) *coalago.CoAPResourceHandlerResult {
	i := info{
		Address: string(storage.Get().Address),
		Version: config.Get("info.version"),
		Type:    config.Get("info.type"),
		Elector: config.Get("everscale.elector"),
		Node:    string(storage.Get().Node),
	}
	info, err := json.Marshal(i)
	if err != nil {
		log.Errorw(err.Error(), "info", i)
		return coalago.NewResponse(coalago.NewStringPayload(err.Error()), coalago.CoapCodeBadRequest)
	}
	log.Debug(string(info))
	return coalago.NewResponse(coalago.NewBytesPayload(info), coalago.CoapCodeContent)
}

func ExecCmd(message *coalago.CoAPMessage) *coalago.CoAPResourceHandlerResult {
	log.Debug(message.Payload.String())
	// decrypt

	// parsing message from node
	command := dsm.CMD{}
	err := json.Unmarshal(message.Payload.Bytes(), &command)
	if err != nil {
		log.Error(err)
		return coalago.NewResponse(coalago.NewStringPayload(err.Error()), coalago.CoapCodeBadRequest)
	}

	// check if the command is sent by one of the device owners
	if !storage.IsOwner(command.Sender) {
		return coalago.NewResponse(coalago.NewStringPayload("invalid sender"), coalago.CoapCodeUnauthorized)
	}

	// verify signature
	// for production: only valid signature is allowed
	// for other env: "testing" can be used as a signature
	if !command.Valid() {
		return coalago.NewResponse(coalago.NewStringPayload("invalid cmd"), coalago.CoapCodeUnauthorized)
	}
	hash, valid := crypto.KeyPair.Verify(command.Sign)
	if (!valid || hash != command.GetHash()) && (os.Args[1] == "prod" || command.Sign != "testing") {
		return coalago.NewResponse(coalago.NewStringPayload("invalid signature"), coalago.CoapCodeUnauthorized)
	}

	// exec command from node
	log.Debug(command.Body)
	cmdArr := strings.Split(command.Body, " ")
	var args []string
	if len(cmdArr) > 1 {
		args = cmdArr[1:]
	}
	log.Debug(cmdArr[0], args)
	c := exec.Command(cmdArr[0], args...)
	if errors.Is(c.Err, exec.ErrDot) {
		c.Err = nil
	}
	log.Debug(c.String(), args)

	stderr, _ := c.StderrPipe()
	stdout, _ := c.StdoutPipe()
	if err = c.Start(); err != nil {
		log.Error(err)
		return coalago.NewResponse(coalago.NewStringPayload(err.Error()), coalago.CoapCodeInternalServerError)
	}

	var errOut string
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		errOut += scanner.Text() + "\n"
	}
	if len(errOut) > 0 {
		log.Error(errOut)
		return coalago.NewResponse(coalago.NewStringPayload(errOut), coalago.CoapCodeInternalServerError)
	}

	var out string
	scanner = bufio.NewScanner(stdout)
	for scanner.Scan() {
		out += scanner.Text() + "\n"
	}
	log.Debug(out)
	return coalago.NewResponse(coalago.NewStringPayload(out), coalago.CoapCodeContent)
}
