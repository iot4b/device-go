package handlers

import (
	"device-go/packages/cmd"
	"device-go/packages/crypto"
	"device-go/packages/dsm"
	"device-go/packages/storage"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/jinzhu/copier"

	"github.com/coalalib/coalago"
	"github.com/google/uuid"
	log "github.com/ndmsystems/golog"
)

type info struct {
	Address   dsm.EverAddress `json:"address"`
	PublicKey string          `json:"publicKey"`
}

// GetInfo info для коалы
func GetInfo(_ *coalago.CoAPMessage) *coalago.CoAPResourceHandlerResult {
	i := info{
		Address:   storage.Device.Address,
		PublicKey: crypto.Keys.PublicSign,
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
	result, err := command.Execute()
	if err != nil {
		return coalago.NewResponse(coalago.NewStringPayload("command.Execute :"+err.Error()), coalago.CoapCodeInternalServerError)
	}

	// generate response with signature
	out := cmd.CMD{
		UUID:   uuid.New().String(),
		Ts:     time.Now().Unix(),
		Sender: crypto.Keys.PublicSign,
		Body:   result,
	}
	hash := out.GetHash()
	out.Hash = hex.EncodeToString(hash)

	sign := crypto.Keys.Sign(hash)
	out.Sign = base64.StdEncoding.EncodeToString(sign)

	s, _ := json.Marshal(out)
	return coalago.NewResponse(coalago.NewStringPayload(string(s)), coalago.CoapCodeContent)
}

// Update local device info with actual data from blockchain
func Update(message *coalago.CoAPMessage) *coalago.CoAPResourceHandlerResult {
	log.Debug("Update:", message.Payload.String())

	var payload struct {
		Address          dsm.EverAddress            `json:"address"`
		Name             string                     `json:"name"`
		Group            dsm.EverAddress            `json:"group"`
		Node             dsm.EverAddress            `json:"node"`
		DeviceAPI        dsm.EverAddress            `json:"deviceAPI"`
		Owners           map[string]dsm.EverAddress `json:"owners"`
		Lock             bool                       `json:"lock,omitempty"`
		Stat             bool                       `json:"stat,omitempty"`
		Events           bool                       `json:"events,omitempty"`
		Version          string                     `json:"version,omitempty"`
		LastRegisterTime string                     `json:"lastRegisterTime,omitempty"`
		NodePubKey       string                     `json:"nodePubKey"`
		Signature        string                     `json:"signature"`
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
	if !crypto.VerifySignature(payload.NodePubKey, []byte(cur), payload.Signature) {
		prev := now.Add(-time.Minute).Format(format)
		if !crypto.VerifySignature(payload.NodePubKey, []byte(prev), payload.Signature) {
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
	signatureHex := hex.EncodeToString(signature)

	return coalago.NewResponse(coalago.NewStringPayload(signatureHex), coalago.CoapCodeContent)
}
