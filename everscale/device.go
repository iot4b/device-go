package everscale

import (
	"device-go/crypto"
	"device-go/dsm"
	log "github.com/ndmsystems/golog"
	"github.com/pkg/errors"
)

var Device = device{}

type device struct {
	Address dsm.EverAddress
}

// SetNode to device smartcontract
func (d device) SetNode(node dsm.EverAddress) error {
	input := map[string]any{"value": node}
	s := NewSigner(crypto.Keys.PublicSign, crypto.Keys.Secret)
	r, err := execute("Device", d.Address, "setNode", input, s)
	if err != nil {
		log.Debug(err, "setNode: "+string(r), input)
		return errors.Wrap(err, "setNode")
	}
	return nil
}
