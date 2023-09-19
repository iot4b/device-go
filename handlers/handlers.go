package handlers

import (
	"bufio"
	"device-go/dsm"
	"device-go/shared"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
	"time"

	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
)

// TOdO везде добавить описания методов и полей моделей

// GetInfo - получить информацию о девайсе
func GetInfo(_ *coalago.CoAPMessage) *coalago.CoAPResourceHandlerResult {
	shared.Info.Uptime = time.Since(shared.Info.RunFrom).String()

	result, err := json.Marshal(shared.Info)
	if err != nil {
		log.Error(err)
		return nil
	}

	log.Debug("device info", shared.Info)

	handlerResult := coalago.NewResponse(coalago.NewStringPayload(string(result)), coalago.CoapCodeContent)
	log.Debug(handlerResult)
	return handlerResult
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

	// todo разобрать команду приватным ключом из sight, иначе вернуть ошибку
	// exec command from node
	log.Debug(command.Cmd)
	cmdArr := strings.Split(command.Cmd, " ")
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
