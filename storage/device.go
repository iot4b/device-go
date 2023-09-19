package storage

import (
	"device-go/cryptoKeys"
	"device-go/dsm"
	"device-go/utils"
	"encoding/json"
	"github.com/jinzhu/copier"
	log "github.com/ndmsystems/golog"
	"github.com/pkg/errors"
	"os"
)

var (
	currentDevice *dsm.DeviceContract
	localPath     string
)

// todo взаимодействие с Elector'ом

func Init(path, vendor, Type, version, vendorData string) {
	currentDevice = new(dsm.DeviceContract)
	localPath = path

	log.Info("Init Local Storage")

	localData, err := readFromLocalStorage(localPath)
	if err != nil && (errors.Is(err, utils.ErrUnmarshal) || errors.Is(err, os.ErrNotExist)) {
		// todo заменить mock данные на реальные адреса в блокчейне
		// локальный файл не найден, инициируем пустой mock
		device := dsm.DeviceContract{
			Address:       "",
			Node:          "",
			Elector:       "",
			VendorAddress: utils.GenerateRandomAddress(),

			PublicKey: cryptoKeys.KeyPair.Public(),
			SecretKey: cryptoKeys.KeyPair.Secret(),

			Stat:       false,
			Type:       Type,
			Version:    version,
			VendorName: vendor,
			VendorData: vendorData,
		}
		err = writeToLocalStorage(localPath, device)
		if err != nil {
			log.Fatal(err)
		}
		Set(device)
		log.Debug("load device contract data from empty obj")
		return
	}
	// если все ок прочиталось из файла, то уст. данные из дампа
	Set(localData)
	log.Debug("load device contract data from local dump file")
}

func Set(d dsm.DeviceContract) {
	copier.Copy(currentDevice, d)
	log.Debugw("local storage device set", "value", *currentDevice)
}

func Update(d *dsm.DeviceContract) error {
	log.Debugw("local storage device update", "before", *currentDevice)

	// не был задан адрес, значит новый девайс
	if len(currentDevice.Address) == 0 && len(d.Address) > 0 {
		currentDevice.Address = d.Address
	}
	// электор задается один раз при старте и в дальнейшем, только vendor может изменить
	if len(currentDevice.Elector) == 0 && len(d.Elector) > 0 {
		currentDevice.Elector = d.Elector
	}
	// Node всегда перезаписывается
	currentDevice.Node = d.Node
	currentDevice.Stat = d.Stat

	log.Debugw("Local Storage Device update", "after", *currentDevice)

	// сохраняем в файл
	return writeToLocalStorage(localPath, *currentDevice)
}

func Get() *dsm.DeviceContract {
	return currentDevice
}

// writeToLocalStorage - сохраняем ноду локально
func writeToLocalStorage(path string, d dsm.DeviceContract) error {
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
