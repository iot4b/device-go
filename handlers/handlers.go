package handlers

import (
	"device-go/aliver"
	"device-go/cmd"
	"device-go/crypto"
	"device-go/everscale"
	"device-go/registration"
	"device-go/shared/config"
	"device-go/storage"
	"encoding/json"
	"time"

	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
)

type info struct {
	Address    string `json:"address"`
	Version    string `json:"version"`
	Elector    string `json:"elector"`
	Node       string `json:"node"`
	Type       string `json:"type"`
	PublicSign string `json:"public_sign"`
	PublicNacl string `json:"public_nacl"`
}

// info для коалы
func GetInfo(_ *coalago.CoAPMessage) *coalago.CoAPResourceHandlerResult {
	i := info{
		Address:    string(storage.Get().Address),
		Version:    config.Get("info.version"),
		Type:       config.Get("info.type"),
		Elector:    config.Get("everscale.elector"),
		Node:       string(storage.Get().Node),
		PublicSign: crypto.Keys.PublicSign,
		PublicNacl: crypto.Keys.PublicNacl,
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

	// validate command params
	if !command.Valid() {
		return coalago.NewResponse(coalago.NewStringPayload("invalid cmd params"), coalago.CoapCodeUnauthorized)
	}
	// verify signature
	if !command.VerifySignature() {
		return coalago.NewResponse(coalago.NewStringPayload("invalid signature"), coalago.CoapCodeUnauthorized)
	}
	// execute command
	out, err := command.Execute()
	if err != nil {
		return coalago.NewResponse(coalago.NewStringPayload(err.Error()), coalago.CoapCodeBadRequest)
	}
	return coalago.NewResponse(coalago.NewStringPayload(out), coalago.CoapCodeContent)
}

// Confirm registration on the node
func Confirm(message *coalago.CoAPMessage) *coalago.CoAPResourceHandlerResult {
	log.Debug(message.Payload.String())

	// update registration data
	_, nodeHost, err := registration.Register()
	if err != nil {
		return coalago.NewResponse(coalago.NewStringPayload(err.Error()), coalago.CoapCodeInternalServerError)
	}
	aliver.NodeHost = nodeHost

	return coalago.NewResponse(coalago.NewStringPayload(""), coalago.CoapCodeContent)
}

// Update local device info with actual data from blockchain
func Update(message *coalago.CoAPMessage) *coalago.CoAPResourceHandlerResult {
	log.Debug(message.Payload.String())

	var payload struct {
		Sign   string `json:"sign"`
		Pubkey string `json:"pubkey"`
	}
	err := json.Unmarshal(message.Payload.Bytes(), &payload)
	if err != nil {
		log.Error(err)
		return coalago.NewResponse(coalago.NewStringPayload(err.Error()), coalago.CoapCodeBadRequest)
	}

	// verify signature
	format := "2006-01-02 15:04"
	now := time.Now()
	cur := now.Format(format)
	if !crypto.Keys.VerifySignature(payload.Pubkey, []byte(cur), payload.Sign) {
		prev := now.Add(-time.Minute).Format(format)
		if !crypto.Keys.VerifySignature(payload.Pubkey, []byte(prev), payload.Sign) {
			return coalago.NewResponse(coalago.NewStringPayload("invalid signature"), coalago.CoapCodeBadRequest)
		}
	}

	// get data from device contract
	d, err := everscale.Device.Get()
	if err != nil {
		return coalago.NewResponse(coalago.NewStringPayload(err.Error()), coalago.CoapCodeInternalServerError)
	}
	// write to local file
	err = storage.WriteToLocalStorage(config.Get("localFiles.contract"), d)
	if err != nil {
		return coalago.NewResponse(coalago.NewStringPayload(err.Error()), coalago.CoapCodeInternalServerError)
	}
	// set data to memory
	storage.Set(d)

	return coalago.NewResponse(coalago.NewStringPayload(""), coalago.CoapCodeContent)
}
