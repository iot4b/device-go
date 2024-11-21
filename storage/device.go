package storage

import (
	"device-go/dsm"
	"device-go/shared/config"
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

func Init(path, elector, vendor, vendorName, vendorData, Type, version string, owners map[string]any) {
	localPath = path

	log.Info("Init Local Storage")
	log.Debug(path, elector, vendor, vendorName, vendorData, Type, version, owners)

	var err error

	// чекаем локально наличие файла
	if utils.FileExists(localPath) {
		Device, err = read(localPath)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		group, err := getGroupAddress()
		if err != nil {
			log.Fatal("failed to read device group address file")
		}
		log.Debug("Group Address:", group)

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
	err = utils.SaveFile(localPath, data)
	if err != nil {
		return errors.Wrapf(err, "utils.SaveFile(%s, data)", localPath)
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

// getGroupAddress get actual device group address from file in config.localFiles.groupAddr
func getGroupAddress() (string, error) {
	addr, err := utils.ReadFile(config.Get("localFiles.groupAddr"))
	if err != nil {
		return "", err
	}
	return string(addr), err
}
