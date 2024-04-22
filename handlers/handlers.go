package handlers

import (
	"device-go/cmd"
	"device-go/crypto"
	"device-go/dsm"
	"device-go/shared/config"
	"device-go/storage"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/jinzhu/copier"

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
	if storage.Get().Lock {
		return coalago.NewResponse(coalago.NewStringPayload("device is locked"), coalago.CoapCodeForbidden)
	}
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
		return coalago.NewResponse(coalago.NewStringPayload(out+" err:"+err.Error()), coalago.CoapCodeBadRequest)
	}
	return coalago.NewResponse(coalago.NewStringPayload(out), coalago.CoapCodeContent)
}

// Update local device info with actual data from blockchain
func Update(message *coalago.CoAPMessage) *coalago.CoAPResourceHandlerResult {
	log.Debug(message.Payload.String())

	var payload struct {
		Address          dsm.EverAddress `json:"address"`
		Node             dsm.EverAddress `json:"node"`
		Active           bool            `json:"active,omitempty"`
		Lock             bool            `json:"lock,omitempty"`
		Stat             bool            `json:"stat,omitempty"`
		Version          string          `json:"version,omitempty"`
		LastRegisterTime string          `json:"lastRegisterTime,omitempty"`
		NodePubKey       string          `json:"nodePubKey"`
		Signature        string          `json:"signature"`
	}
	err := json.Unmarshal(message.Payload.Bytes(), &payload)
	if err != nil {
		log.Error(err)
		return coalago.NewResponse(coalago.NewStringPayload(err.Error()), coalago.CoapCodeBadRequest)
	}

	device := storage.Get()
	if payload.Address != device.Address {
		return coalago.NewResponse(coalago.NewStringPayload("wrong device address"), coalago.CoapCodeBadRequest)
	}

	// verify signature
	format := "2006-01-02 15:04"
	now := time.Now()
	cur := now.Format(format)
	if !crypto.Keys.VerifySignature(payload.NodePubKey, []byte(cur), payload.Signature) {
		prev := now.Add(-time.Minute).Format(format)
		if !crypto.Keys.VerifySignature(payload.NodePubKey, []byte(prev), payload.Signature) {
			return coalago.NewResponse(coalago.NewStringPayload("invalid signature"), coalago.CoapCodeBadRequest)
		}
	}

	// update local data from payload
	copier.Copy(device, payload)

	// write to local file
	err = storage.WriteToLocalStorage(config.Get("localFiles.contract"), *device)
	if err != nil {
		return coalago.NewResponse(coalago.NewStringPayload(err.Error()), coalago.CoapCodeInternalServerError)
	}
	// set data to memory
	storage.Set(*device)

	return coalago.NewResponse(coalago.NewStringPayload(""), coalago.CoapCodeContent)
}

// Sign data with key pair and return signature
func Sign(message *coalago.CoAPMessage) *coalago.CoAPResourceHandlerResult {
	log.Debug(message.Payload.String())

	var payload struct {
		Unsigned string `json:"u"`
	}
	err := json.Unmarshal(message.Payload.Bytes(), &payload)
	if err != nil {
		log.Error(err)
		return coalago.NewResponse(coalago.NewStringPayload(err.Error()), coalago.CoapCodeBadRequest)
	}

	unsigned, err := base64.StdEncoding.DecodeString(payload.Unsigned)
	if err != nil {
		log.Error(err)
		return coalago.NewResponse(coalago.NewStringPayload(err.Error()), coalago.CoapCodeBadRequest)
	}

	signature := crypto.Keys.Sign(unsigned)

	return coalago.NewResponse(coalago.NewStringPayload(signature), coalago.CoapCodeContent)
}
