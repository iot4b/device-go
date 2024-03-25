package storage

import (
	"device-go/dsm"
	"device-go/utils"
	"encoding/json"
	"github.com/jinzhu/copier"
	log "github.com/ndmsystems/golog"
	"github.com/pkg/errors"
)

var (
	currentDevice *dsm.DeviceContract
	localPath     string
)

type initialData struct {
	Elector dsm.EverAddress `json:"elector"`
	Vendor  dsm.EverAddress `json:"vendor"`
	Owners  []string        `json:"owners"`

	Type       string `json:"dtype"`
	Version    string `json:"version"`
	VendorName string `json:"vendorName"`
	VendorData string `json:"vendorData"`
}

var (
	ErrInvalidValue = errors.New("invalid value")
	ErrIsRequired   = errors.New("field is required")
	ErrIsEmpty      = errors.New("value is empty")
	ErrNotSpecified = errors.New("value not specified")
)

// todo взаимодействие с Elector'ом

func Init(path, elector, vendor, vendorName, vendorData, Type, version string, owners []string) {
	currentDevice = new(dsm.DeviceContract)
	localPath = path

	log.Info("Init Local Storage")
	log.Debug(path, elector, vendor, vendorName, vendorData, Type, version, owners)

	var localData dsm.DeviceContract
	var err error

	// чекаем локально наличие файла
	if utils.FileExists(localPath) {
		localData, err = readFromLocalStorage(localPath)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		localData = dsm.DeviceContract{
			Elector:    dsm.EverAddress(elector),
			Vendor:     dsm.EverAddress(vendor),
			Owners:     owners,
			Type:       Type,
			Version:    version,
			VendorName: vendorName,
			VendorData: vendorData,
		}
	}

	Set(localData)
}

func Set(d dsm.DeviceContract) {
	copier.Copy(currentDevice, d)
	log.Debugw("local storage device set", "value", *currentDevice)
}

func Get() *dsm.DeviceContract {
	return currentDevice
}

// WriteToLocalStorage - сохраняем ноду локально
func WriteToLocalStorage(path string, d dsm.DeviceContract) error {
	data, err := json.Marshal(d)
	if err != nil {
		return errors.Wrap(err, "json.Marshal(device)")
	}
	err = utils.SaveFile(path, data)
	if err != nil {
		return errors.Wrapf(err, "utils.SaveFile(%s, data)", path)
	}
	return nil
}

// readFromLocalStorage - читаем ноду из локального дампа контракта
func readFromLocalStorage(path string) (d dsm.DeviceContract, err error) {
	err = utils.ReadJSONFile(path, &d)
	if err != nil {
		return dsm.DeviceContract{}, err
	}
	log.Debugf("%+v", d)
	return d, err
}

func (d initialData) validate() error {
	if len(d.Vendor) == 0 {
		return errors.Wrap(ErrIsRequired, "vendor")
	}
	if len(d.Owners) == 0 {
		return errors.Wrap(ErrIsEmpty, "owners")
	}
	if len(d.Owners) > 0 {
		// todo сделать валидатор ever адресов
		for i, owner := range d.Owners {
			if len(owner) == 0 {
				return errors.Wrapf(ErrInvalidValue, "owners[%d]", i)
			}
		}
	}
	if len(d.VendorName) == 0 {
		return errors.Wrap(ErrNotSpecified, "vendorName")
	}
	if len(d.Type) == 0 {
		return errors.Wrap(ErrNotSpecified, "type")
	}
	if len(d.Version) == 0 {
		return errors.Wrap(ErrNotSpecified, "version")
	}
	return nil
}

func (d initialData) toMap() (result map[string]interface{}) {
	data, _ := json.Marshal(d)
	json.Unmarshal(data, &result)
	return
}

// IsOwner checks if key is one of the owners from device contract
func IsOwner(key string) bool {
	for _, owner := range Get().Owners {
		if key == owner {
			return true
		}
	}
	return false
}
