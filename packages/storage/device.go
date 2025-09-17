package storage

import (
	"device-go/packages/dsm"
	"device-go/packages/utils"
	"encoding/json"
	"os"
	"path/filepath"

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

	Owners map[string]any `json:"owners"` // owners data: public_key => contract_address

	Lock    bool   `json:"lock"`              // if device is locked
	Stat    bool   `json:"stat"`              // нужно ли девайсу слать статистику
	Events  bool   `json:"events"`            // sending events
	Type    string `json:"dtype,omitempty"`   // модель/тип девайса
	Version string `json:"version,omitempty"` // версия текущей прошивки на девайсе

	LastRegisterTime string `json:"lastRegisterTime,omitempty"` // last registration timestamp

	Hash string `json:"hash"` // hash of current contract code (contract version identifier)
}

var (
	Device   device
	filePath string
)

func Init(path, initFile, elector, vendor, deviceAPI, dType, version, group, owner string) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	filePath = filepath.Join(home, ".config", utils.FilesDir, path)

	log.Info("Init Local Storage")
	log.Debug(path, elector, vendor, deviceAPI, dType, version)

	// чекаем локально наличие файла
	if utils.FileExists(filePath) {
		Device, err = read(filePath)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// read from init params from file in config.localFiles.init
		Device, err = read(initFile)
		if err != nil {
			if group == "" || owner == "" {
				log.Fatal("group and owner is required")
			}
			name, err := os.Hostname()
			if err != nil {
				name = "iot4b-device"
			}
			Device = device{
				Name:   name,
				Group:  dsm.EverAddress(group),
				Owners: map[string]any{owner: "0:0000000000000000000000000000000000000000000000000000000000000000"},
			}
		}
		Device.Elector = dsm.EverAddress(elector)
		Device.Vendor = dsm.EverAddress(vendor)
		Device.DeviceAPI = dsm.EverAddress(deviceAPI)
		Device.Type = dType
		Device.Version = version

		if err = Save(); err != nil {
			log.Errorf("storage.Save: %v", err)
		}
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

// IsOwner checks if key is one of the owners from device contract
func IsOwner(key string) bool {
	for owner := range Device.Owners {
		if "0x"+key == owner {
			return true
		}
	}
	return false
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
