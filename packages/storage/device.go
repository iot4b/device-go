package storage

import (
	"device-go/packages/dsm"
	"device-go/packages/utils"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	log "github.com/ndmsystems/golog"
	"github.com/pkg/errors"
)

type device struct {
	Name      string          `json:"name,omitempty"`    //Device name
	Address   dsm.EverAddress `json:"address,omitempty"` //ever SC address текущего Device
	Group     dsm.EverAddress `json:"group,omitempty"`   //ever SC address of DeviceGroup
	Node      dsm.EverAddress `json:"node,omitempty"`    //ever SC address Node, с которой девайс создал последнее соединение
	Elector   dsm.EverAddress `json:"elector"`           //ever SC адрес Elector'a, который обслуживает сеть нод для текущего девайса
	Vendor    dsm.EverAddress `json:"vendor"`            //ever SC address производителя текущего девайса. по-умолчанию из конфигов берем
	DeviceAPI dsm.EverAddress `json:"deviceAPI"`         //ever SC address of Device API contract

	Owners []string `json:"owners"` // owner public keys

	Lock    bool   `json:"lock"`              // if device is locked
	Stat    bool   `json:"stat"`              // нужно ли девайсу слать статистику
	Events  bool   `json:"events"`            // sending events
	Type    string `json:"dtype,omitempty"`   // модель/тип девайса
	Version string `json:"version,omitempty"` // версия текущей прошивки на девайсе

	LastRegisterTime string `json:"lastRegisterTime,omitempty"` // last registration timestamp

	Hash string `json:"hash"` // hash of current contract code (contract version identifier)

	NodeIpPort string `json:"nodeIpPort,omitempty"`
}

var (
	Device   device
	filePath string
)

func Init(path, setupPath, elector, vendor, deviceAPI, dType, deviceVersion string) {
	filePath = filepath.Join(utils.GetFilesDir(), path)
	setupFilePath = resolveSetupFilePath(setupPath)

	log.Debug(path, elector, vendor, deviceAPI, dType, deviceVersion)

	var err error
	// чекаем локально наличие файла
	if utils.FileExists(filePath) {
		Device, err = read(filePath)
		if err != nil {
			log.Fatal(err)
		}
	} else if setup, setupErr := readSetup(); setupErr == nil && setup.Address != "" {
		Device = newInitialDevice()
		Device.Address = dsm.EverAddress(setup.Address)
		Pairing = setup.Pairing
	} else {
		Device = newInitialDevice()
	}
	if setup, setupErr := readSetup(); setupErr == nil {
		Pairing = setup.Pairing
	}

	if Device.Name == "" {
		if host, hostErr := os.Hostname(); hostErr == nil && host != "" {
			Device.Name = host
		}
	}
	if Device.Owners == nil {
		Device.Owners = []string{}
	}
	Device.Elector = dsm.EverAddress(elector)
	Device.Vendor = dsm.EverAddress(vendor)
	Device.DeviceAPI = dsm.EverAddress(deviceAPI)
	Device.Type = dType
	Device.Version = deviceVersion
	Device.NodeIpPort = "" // should be empty before first registration

	if err = Save(); err != nil {
		log.Errorf("storage.Save: %v", err)
	}
}

// Save local data to file
func Save() error {
	data, err := json.Marshal(Device)
	if err != nil {
		return errors.Wrap(err, "json.Marshal(device)")
	}
	err = utils.SaveFile(filePath, data)
	if err != nil {
		return errors.Wrapf(err, "utils.SaveFile(%s, data)", filePath)
	}
	return nil
}

// Update Device data from file
func Update() (err error) {
	if utils.FileExists(filePath) {
		Device, err = read(filePath)
		if err != nil {
			return err
		}
	}
	if setup, setupErr := readSetup(); setupErr == nil {
		Pairing = setup.Pairing
		if Device.Address == "" && setup.Address != "" {
			Device.Address = dsm.EverAddress(setup.Address)
		}
	}
	return nil
}

// IsOwner checks if key is one of the owners from device contract
func IsOwner(key string) bool {
	normalizedKey := normalizeOwnerKey(key)
	for _, owner := range Device.Owners {
		if normalizeOwnerKey(owner) == normalizedKey {
			return true
		}
	}
	return false
}

func normalizeOwnerKey(key string) string {
	return strings.TrimPrefix(strings.ToLower(strings.TrimSpace(key)), "0x")
}

// read local data from file
func read(path string) (d device, err error) {
	err = utils.ReadJSONFile(path, &d)
	if err != nil {
		return device{}, err
	}
	log.Debugf("%+v", d)
	return d, err
}

func HasContractAddress() bool {
	return Device.Address != ""
}

func newInitialDevice() device {
	name, err := os.Hostname()
	if err != nil {
		name = "iot4b-device"
	}

	return device{
		Name:   name,
		Owners: []string{},
	}
}
