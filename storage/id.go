package storage

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"os"
)

const idFile = "id.json"

var Id = id{}

type id struct {
	Address string `json:"a"`   // deployed contract address
	Public  string `json:"pub"` // keys
	Secret  string `json:"sec"` // for signing
}

func init() {
	file, err := os.Open(idFile)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		data, _ := json.Marshal(Id)
		os.WriteFile(idFile, data, 0644)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return
	}
	json.Unmarshal(data, &Id)
}

func (id *id) Save() {
	data, _ := json.Marshal(id)
	os.WriteFile(idFile, data, 0644)
}
