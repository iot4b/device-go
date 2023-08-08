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
	log.Debug(message)
	command := cmd{}
	// parsing message from node
	log.Debug(message)
	err := json.Unmarshal(message.Payload.Bytes(), &command)
	if err != nil {
		log.Error(err)
		return coalago.NewResponse(coalago.NewStringPayload(err.Error()), coalago.CoapCodeInternalServerError)
	}

	// todo разобрать команду приватным ключом из sight, иначе вернуть ошибку
	// exec command from node
	c := exec.Command(command.Cmd)
	if errors.Is(c.Err, exec.ErrDot) {
		c.Err = nil
	}
	if err := c.Run(); err != nil {
		log.Error(err)
		return coalago.NewResponse(coalago.NewStringPayload(err.Error()), coalago.CoapCodeInternalServerError)
	}

	// todo что отдавать в ответ после выполнения колманды
	return coalago.NewResponse(nil, coalago.CoapCodeEmpty)
}
