package cmd

import "device-go/packages/dsm"

type CMD struct {
	UUID       string          `json:"uuid"`        // unique id
	Ts         int64           `json:"ts"`          // timestamp
	Sender     string          `json:"sender"`      // owner public key
	SenderNacl string          `json:"sender_nacl"` // sender public key for nacl box encryption
	Receiver   dsm.EverAddress `json:"receiver"`    // device contract address
	Hash       string          `json:"hash"`        // hash all fields except Sign and Hash
	Sign       string          `json:"sign"`        // signature of hash with private key of sender
	Body       string          `json:"body"`        // command to execute, encrypted with nacl box, first 48 chars is a nonce
}
