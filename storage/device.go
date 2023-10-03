package storage

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"os"
)

const deviceFile = "device.json"

var Device = device{}

type device struct {
	Address string `json:"a"` // deployed contract address
	Public  string `json:"p"` // keys
	Secret  string `json:"x"` // for signing
}

func init() {
	file, err := os.Open(deviceFile)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		data, _ := json.Marshal(Device)
		os.WriteFile(deviceFile, data, 0644)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return
	}
	json.Unmarshal(data, &Device)
}

func (id *device) Save() {
	data, _ := json.Marshal(id)
	os.WriteFile(deviceFile, data, 0644)
}
