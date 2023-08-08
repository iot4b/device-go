package api

import (
	"encoding/json"
	"errors"
	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
	"os/exec"
)

type cmd struct {
	Cmd string `json:"cmd"`
	// todo для чего sight и cmd раздельно?
	Sight string `json:"sight"`
	Uid   string `json:"uid"`
}

func getInfo(_ *coalago.CoAPMessage) *coalago.CoAPResourceHandlerResult {
	result, err := json.Marshal(info)
	if err != nil {
		log.Error(err)
		return nil
	}

	handlerResult := coalago.NewResponse(coalago.NewStringPayload(string(result)), coalago.CoapCodeContent)
	log.Debug(handlerResult)
	return handlerResult
}

func execCmd(message *coalago.CoAPMessage) *coalago.CoAPResourceHandlerResult {
	log.Debug(message.Payload.String())
	// parsing message from node
	command := cmd{}
	err := json.Unmarshal(message.Payload.Bytes(), &command)
	if err != nil {
		log.Error(err)
		return coalago.NewResponse(coalago.NewStringPayload(err.Error()), coalago.CoapCodeBadRequest)
	}

	// todo разобрать команду приватным ключом из sight, иначе вернуть ошибку
	// exec command from node
	c := exec.Command(command.Cmd)
	if errors.Is(c.Err, exec.ErrDot) {
		c.Err = nil
	}
	output, err := c.Output()
	if err != nil {
		log.Error(err)
		return coalago.NewResponse(coalago.NewStringPayload(err.Error()), coalago.CoapCodeInternalServerError)
	}

	log.Debug(string(output))
	return coalago.NewResponse(coalago.NewStringPayload(string(output)), coalago.CoapCodeContent)
}
