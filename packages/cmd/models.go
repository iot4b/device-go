package cmd

import "device-go/packages/dsm"

type CMD struct {
	UUID     string          `json:"uuid"`     // unique id
	Ts       int64           `json:"ts"`       // timestamp
	Sender   string          `json:"sender"`   // owner public key
	Receiver dsm.EverAddress `json:"receiver"` // device contract address
	Hash     string          `json:"hash"`     // hash all fields except Sign and Hash
	Sign     string          `json:"sign"`     // signature of hash with private key of sender
	Body     string          `json:"body"`     // command to execute, base64-encoded, encrypted with chacha20-poly1305, first 12 bytes of data is a nonce
}
