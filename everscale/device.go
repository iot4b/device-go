package everscale

import (
	"device-go/crypto"
	"device-go/dsm"
	"encoding/json"
)

var Device = device{}

type device struct {
	Address dsm.EverAddress
}

// SetNode to device smartcontract
func (d device) SetNode(node dsm.EverAddress) error {
	input := map[string]any{"value": node}
	s := newSigner(crypto.Keys.PublicSign, crypto.Keys.Secret)
	_, err := execute("Device", d.Address, "setNode", input, s)
	return err
}

// Get actual device data from blockchain
func (d device) Get() (device dsm.DeviceContract, err error) {
	s := newSigner(crypto.Keys.PublicSign, crypto.Keys.Secret)
	r, err := execute("Device", d.Address, "get", nil, s)
	if err != nil {
		return
	}
	err = json.Unmarshal(r, &device)
	if err != nil {
		return
	}
	device.Address = d.Address
	return
}
