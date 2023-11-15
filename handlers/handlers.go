package handlers

import (
	"device-go/cmd"
	"device-go/shared/config"
	"device-go/storage"
	"encoding/json"

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
	// parsing message from node
	command, err := cmd.Build(message.Payload.Bytes())
	if err != nil {
		log.Error(err)
		return coalago.NewResponse(coalago.NewStringPayload(err.Error()), coalago.CoapCodeBadRequest)
	}
	log.Debug(command.Readable())

	// check if the command is sent by one of the device owners
	if !storage.IsOwner(command.Sender) {
		return coalago.NewResponse(coalago.NewStringPayload("invalid sender"), coalago.CoapCodeUnauthorized)
	}

	// verify signature
	// for production: only valid signature is allowed
	// for other env: "testing" can be used as a signature
	if !command.Valid() {
		return coalago.NewResponse(coalago.NewStringPayload("invalid cmd params"), coalago.CoapCodeUnauthorized)
	}
	// verify signature
	if _, result := command.Verify(); !result {
		return coalago.NewResponse(coalago.NewStringPayload("invalid signature"), coalago.CoapCodeUnauthorized)
	}
	// execute command
	out, err := command.Execute()
	if err != nil {
		return coalago.NewResponse(coalago.NewStringPayload(err.Error()), coalago.CoapCodeBadRequest)
	}
	return coalago.NewResponse(coalago.NewStringPayload(out), coalago.CoapCodeContent)
}
