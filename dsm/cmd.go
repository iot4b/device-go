package dsm

import (
	"crypto/sha256"
)

type CMD struct {
	UUID     string      `json:"uuid"`     // unique id
	Ts       int64       `json:"ts"`       // timestamp
	Sender   string      `json:"sender"`   // owner public key
	Receiver EverAddress `json:"receiver"` // device contract address
	Hash     string      `json:"hash"`     // hash all fields except Sign and Hash
	Sign     string      `json:"sign"`     // signature of hash with private key of sender
	Body     string      `json:"body"`     // command to execute
}

// Valid checks if all fields are filled
func (c CMD) Valid() bool {
	if len(c.UUID) == 0 {
		return false
	}
	if c.Ts == 0 {
		return false
	}
	if len(c.Sender) == 0 {
		return false
	}
	if len(c.Receiver) == 0 {
		return false
	}
	if len(c.Hash) == 0 {
		return false
	}
	if len(c.Sign) == 0 {
		return false
	}
	return true
}

// GetHash calculates hash sum of all fields except Sign and Hash
func (c CMD) GetHash() string {
	h := sha256.New()
	bt := []byte(string(c.Sender + string(c.Receiver) + c.Body + string(c.Ts) + string(c.UUID)))
	h.Write(bt)
	return string(h.Sum(nil))
}
