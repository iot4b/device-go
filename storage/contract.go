package storage

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"os"
)

const contractFile = "contract.json"

var Contract = contract{}

type contract struct {
	Address string `json:"a"` // deployed contract address
	Public  string `json:"p"` // keys
	Secret  string `json:"x"` // for signing
}

func init() {
	file, err := os.Open(contractFile)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		data, _ := json.Marshal(Contract)
		os.WriteFile(contractFile, data, 0644)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return
	}
	json.Unmarshal(data, &Contract)
}

func (id *contract) Save() {
	data, _ := json.Marshal(id)
	os.WriteFile(contractFile, data, 0644)
}
