package dsm

import (
	"crypto/sha256"
)

type CMD struct {
	Sender   string      `json:"sender"`   // owner public key
	Receiver EverAddress `json:"receiver"` // device contract address
	Hash     string      `json:"hash"`     // hash all fields except Sign and Hash
	Sign     string      `json:"sign"`     // signature of hash with private key of sender
	Body     string      `json:"body"`     // command to execute
}

// GetHash calculates hash sum of all fields except Sign and Hash
func (c CMD) GetHash() string {
	h := sha256.New()
	h.Write([]byte(c.Sender + string(c.Receiver) + c.Body))
	return string(h.Sum(nil))
}
