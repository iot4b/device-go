package storage

import (
	"device-go/packages/utils"
	"encoding/json"
	"path/filepath"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type pairingState struct {
	Code        string `json:"code,omitempty"`
	NodeIpPort  string `json:"nodeIpPort,omitempty"`
	NodeAddress string `json:"nodeAddress,omitempty"`
	Status      string `json:"status,omitempty"`
	HeartbeatAt int64  `json:"heartbeatAt,omitempty"`
}

type setupData struct {
	Address string       `json:"address,omitempty"`
	Pairing pairingState `json:"pairing,omitempty"`
}

var (
	Pairing       pairingState
	setupFilePath string
	setupMu       sync.Mutex
)

func saveSetup() error {
	if setupFilePath == "" {
		return nil
	}

	data, err := json.Marshal(setupData{
		Address: Device.Address.String(),
		Pairing: Pairing,
	})
	if err != nil {
		return errors.Wrap(err, "json.Marshal(setup)")
	}

	return utils.SaveFile(setupFilePath, data)
}

func readSetup() (setupData, error) {
	if setupFilePath == "" || !utils.FileExists(setupFilePath) {
		return setupData{}, errors.New("setup file is not available")
	}

	var setup setupData
	if err := utils.ReadJSONFile(setupFilePath, &setup); err != nil {
		return setupData{}, err
	}

	return setup, nil
}

func resolveSetupFilePath(path string) string {
	if path == "" {
		return ""
	}
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(utils.GetFilesDir(), path)
}

func HasPairingCode() bool {
	return Pairing.Code != ""
}

func SetPairing(code, nodeIpPort, nodeAddress, status string) error {
	setupMu.Lock()
	defer setupMu.Unlock()

	Pairing.Code = code
	Pairing.NodeIpPort = nodeIpPort
	Pairing.NodeAddress = nodeAddress
	Pairing.Status = status
	return saveSetup()
}

func ClearPairing() error {
	setupMu.Lock()
	defer setupMu.Unlock()

	heartbeatAt := Pairing.HeartbeatAt
	Pairing = pairingState{HeartbeatAt: heartbeatAt}
	return saveSetup()
}

func PairingTouch() error {
	setupMu.Lock()
	defer setupMu.Unlock()

	heartbeatAt := time.Now().Unix()
	if Pairing.HeartbeatAt == heartbeatAt {
		return nil
	}
	Pairing.HeartbeatAt = heartbeatAt
	return saveSetup()
}
