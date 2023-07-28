package api

import (
	"encoding/json"
	"errors"
	"github.com/coalalib/coalago/message"
	"github.com/coalalib/coalago/resource"
	log "github.com/ndmsystems/golog"
	"os/exec"
)

type cmd struct {
	Cmd string `json:"cmd"`
	// todo для чего sight и cmd раздельно?
	Sight string `json:"sight"`
	Uid   string `json:"uid"`
}

func getInfo(_ *coalaMsg.CoAPMessage) *resource.CoAPResourceHandlerResult {
	result, err := json.Marshal(info)
	if err != nil {
		log.Error(err)
		return nil
	}

	handlerResult := resource.NewResponse(coalaMsg.NewStringPayload(string(result)), coalaMsg.CoapCodeContent)
	log.Debug(handlerResult)
	return handlerResult
}

func execCmd(message *coalaMsg.CoAPMessage) *resource.CoAPResourceHandlerResult {
	log.Debug(message)
	command := cmd{}
	// parsing message from node
	err := json.Unmarshal(message.Payload.Bytes(), &command)
	if err != nil {
		log.Error(err)
		return nil
	}

	// todo разобрать команду приватным ключом из sight, иначе вернуть ошибку
	// exec command from node
	c := exec.Command(command.Cmd)
	if errors.Is(c.Err, exec.ErrDot) {
		c.Err = nil
	}
	if err := c.Run(); err != nil {
		log.Error(err)
	}

	// todo что отдавать в ответ после выполнения колманды
	return resource.NewResponse(nil, coalaMsg.CoapCodeEmpty)
}
