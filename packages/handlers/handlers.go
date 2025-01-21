package handlers

import (
	"device-go/packages/cmd"
	"device-go/packages/crypto"
	"device-go/packages/dsm"
	"device-go/packages/storage"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/jinzhu/copier"

	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
)

type info struct {
	Address    dsm.EverAddress `json:"address"`
	Group      dsm.EverAddress `json:"group"`
	Node       dsm.EverAddress `json:"node"`
	Elector    dsm.EverAddress `json:"elector"`
	Vendor     dsm.EverAddress `json:"vendor"`
	Owners     map[string]any  `json:"owners"`
	Lock       bool            `json:"lock"`
	Stat       bool            `json:"stat"`
	Events     bool            `json:"events"`
	Type       string          `json:"type"`
	Version    string          `json:"version"`
	VendorName string          `json:"vendorName"`
	PublicSign string          `json:"public_sign"`
	PublicNacl string          `json:"public_nacl"`
}

// info для коалы
func GetInfo(_ *coalago.CoAPMessage) *coalago.CoAPResourceHandlerResult {
	i := info{}
	copier.Copy(&i, storage.Device)
	i.PublicSign = crypto.Keys.PublicSign
	i.PublicNacl = crypto.Keys.PublicNacl

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
	if storage.Device.Lock {
		log.Debug("device locked")
		return coalago.NewResponse(coalago.NewStringPayload("device is locked"), coalago.CoapCodeForbidden)
	}
	// parsing message from node
	command, err := cmd.Build(message.Payload.Bytes())
	if err != nil {
		log.Error(err)
		return coalago.NewResponse(coalago.NewStringPayload(err.Error()), coalago.CoapCodeBadRequest)
	}
	// validate command params
	if !command.Valid() {
		return coalago.NewResponse(coalago.NewStringPayload("invalid cmd params"), coalago.CoapCodeUnauthorized)
	}
	// verify signature
	if !command.VerifySignature() {
		return coalago.NewResponse(coalago.NewStringPayload("invalid signature"), coalago.CoapCodeUnauthorized)
	}
	// check if the command is sent by one of the device owners
	if !storage.IsOwner(command.Sender) {
		log.Debug("invalid sender")
		return coalago.NewResponse(coalago.NewStringPayload("invalid sender"), coalago.CoapCodeUnauthorized)
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
	log.Debug("Update:", message.Payload.String())

	var payload struct {
		Address          dsm.EverAddress `json:"address"`
		Group            dsm.EverAddress `json:"group"`
		Node             dsm.EverAddress `json:"node"`
		Lock             bool            `json:"lock,omitempty"`
		Stat             bool            `json:"stat,omitempty"`
		Events           bool            `json:"events,omitempty"`
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

	if payload.Address != storage.Device.Address {
		return coalago.NewResponse(coalago.NewStringPayload("wrong device address"), coalago.CoapCodeBadRequest)
	}

	// verify signature
	format := "2006-01-02 15:04"
	now := time.Now().UTC()
	cur := now.Format(format)
	if !crypto.Keys.VerifySignature(payload.NodePubKey, []byte(cur), payload.Signature) {
		prev := now.Add(-time.Minute).Format(format)
		if !crypto.Keys.VerifySignature(payload.NodePubKey, []byte(prev), payload.Signature) {
			return coalago.NewResponse(coalago.NewStringPayload("invalid signature"), coalago.CoapCodeBadRequest)
		}
	}

	// update local data from payload
	copier.Copy(&storage.Device, payload)

	// write to local file
	err = storage.Save()
	if err != nil {
		return coalago.NewResponse(coalago.NewStringPayload(err.Error()), coalago.CoapCodeInternalServerError)
	}

	return coalago.NewResponse(coalago.NewStringPayload(""), coalago.CoapCodeContent)
}

// Sign : decrypt nacl box, sign data with sign key pair and return signature
func Sign(message *coalago.CoAPMessage) *coalago.CoAPResourceHandlerResult {
	log.Debug("Sign:", message.Payload.String())

	var payload struct {
		NaclBox   string `json:"b"` // encrypted data, base64
		PublicKey string `json:"k"` // node public key for nacl box
	}
	err := json.Unmarshal(message.Payload.Bytes(), &payload)
	if err != nil {
		log.Error(err)
		return coalago.NewResponse(coalago.NewStringPayload(err.Error()), coalago.CoapCodeBadRequest)
	}

	unsigned, err := crypto.Keys.Decrypt(payload.NaclBox, payload.PublicKey)
	if err != nil {
		log.Error(err)
		return coalago.NewResponse(coalago.NewStringPayload(err.Error()), coalago.CoapCodeBadRequest)
	}

	data, err := base64.StdEncoding.DecodeString(unsigned)
	if err != nil {
		log.Error(err)
		return coalago.NewResponse(coalago.NewStringPayload(err.Error()), coalago.CoapCodeBadRequest)
	}
	signature := crypto.Keys.Sign(data)

	return coalago.NewResponse(coalago.NewStringPayload(signature), coalago.CoapCodeContent)
}
