package storage

import (
	"device-go/dsm"
	"device-go/utils"
	"encoding/json"
	log "github.com/ndmsystems/golog"
	"github.com/pkg/errors"
)

type device struct {
	Name    string          `json:"name,omitempty"`    //Device name
	Address dsm.EverAddress `json:"address,omitempty"` //ever SC address текущего Device
	Group   dsm.EverAddress `json:"group,omitempty"`   //ever SC address of DeviceGroup
	Node    dsm.EverAddress `json:"node,omitempty"`    //ever SC address Node, с которой девайс создал последнее соединение
	Elector dsm.EverAddress `json:"elector"`           //ever SC адрес Elector'a, который обслуживает сеть нод для текущего девайса
	Vendor  dsm.EverAddress `json:"vendor"`            //ever SC address производителя текущего девайса. по-умолчанию из конфигов берем

	Owners map[string]any `json:"owners"` // owners data: public_key => contract_address

	Active     bool   `json:"active"`               // if device is active
	Lock       bool   `json:"lock"`                 // if device is locked
	Stat       bool   `json:"stat"`                 // нужно ли девайсу слать статистику
	Events     bool   `json:"events"`               // sending events
	Type       string `json:"dtype,omitempty"`      // модель/тип девайса
	Version    string `json:"version,omitempty"`    // версия текущей прошивки на девайсе
	VendorName string `json:"vendorName,omitempty"` // название производителя
	VendorData string `json:"vendorData,omitempty"` // данные, которые идут от производителя девайса

	LastRegisterTime string `json:"lastRegisterTime,omitempty"` // last registration timestamp

	Hash string `json:"hash"` // hash of current contract code (contract version identifier)
}

var (
	Device    device
	localPath string
)

func Init(path, elector, vendor, vendorName, vendorData, Type, version string, owners map[string]any, group string) {
	localPath = path

	log.Info("Init Local Storage")
	log.Debug(path, elector, vendor, vendorName, vendorData, Type, version, owners, group)

	var err error

	// чекаем локально наличие файла
	if utils.FileExists(localPath) {
		Device, err = read(localPath)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		Device = device{
			Group:      dsm.EverAddress(group),
			Elector:    dsm.EverAddress(elector),
			Vendor:     dsm.EverAddress(vendor),
			Owners:     owners,
			Type:       Type,
			Version:    version,
			VendorName: vendorName,
			VendorData: vendorData,
		}
		if err = Save(); err != nil {
			log.Errorf("WriteToLocalStorage: %v", err)
		}
	}
}

// Save local data to file
func Save() error {
	data, err := json.Marshal(Device)
	if err != nil {
		return errors.Wrap(err, "json.Marshal(device)")
	}
	err = utils.SaveFile(localPath, data)
	if err != nil {
		return errors.Wrapf(err, "utils.SaveFile(%s, data)", localPath)
	}
	return nil
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

// IsOwner checks if key is one of the owners from device contract
func IsOwner(key string) bool {
	for owner := range Device.Owners {
		if "0x"+key == owner {
			return true
		}
	}
	return false
}
