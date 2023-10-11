package storage

import (
	"device-go/cryptoKeys"
	"device-go/dsm"
	"device-go/everscale"
	"device-go/shared/config"
	"device-go/utils"
	"encoding/json"
	"github.com/jinzhu/copier"
	log "github.com/ndmsystems/golog"
	"github.com/pkg/errors"
	"os"
	"time"
)

var (
	currentDevice *dsm.DeviceContract
	localPath     string
)

type initialData struct {
	Elector dsm.EverAddress   `json:"elector"`
	Vendor  dsm.EverAddress   `json:"vendor"`
	Owners  []dsm.EverAddress `json:"owners"`

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

	// чекаем локально наличие файла
	localData, err := readFromLocalStorage(localPath)
	if err != nil {
		// если файла нет или ошибка формата данных, то деплоим контракт в блокчейн
		if errors.Is(err, utils.ErrUnmarshal) || errors.Is(err, os.ErrNotExist) {
			// todo заменить mock данные на реальные адреса в блокчейне
			// локальный файл не найден, инициируем пустой mock
			o := make([]dsm.EverAddress, 0)
			for _, owner := range owners {
				o = append(o, dsm.EverAddress(owner))
			}
			data := initialData{
				Elector:    dsm.EverAddress(elector),
				Vendor:     dsm.EverAddress(vendor),
				Owners:     o,
				Type:       Type,
				Version:    version,
				VendorName: vendorName,
				VendorData: vendorData,
			}
			log.Debug("initial contract data", data)

			device, err := deploy(cryptoKeys.KeyPair.PublicStr(), cryptoKeys.KeyPair.SecretStr(), data)
			if err != nil {
				log.Fatal(err)
			}
			// сохраняем локально
			err = writeToLocalStorage(localPath, device)
			if err != nil {
				log.Fatal(err)
			}
			Set(device)
			log.Debug("load device contract data from empty obj")
			return
		}
		log.Fatal(err)
	}
	// если все ок прочиталось из файла, то уст. данные из дампа
	Set(localData)
	log.Debug("load device contract data from local dump file")
}

func deploy(public, secret string, data initialData) (out dsm.DeviceContract, err error) {
	log.Info("Deploy device contract")
	// валидируем
	err = data.validate()
	if err != nil {
		return
	}
	log.Debug("validate initial data OK!")

	// giver - это такой кошелек, который по
	abi, tvc, err := everscale.ReadContract("./contracts", "device")
	if err != nil {
		return
	}

	// init ContractBuilder
	device := &everscale.ContractBuilder{Public: public, Secret: secret, Abi: abi, Tvc: tvc}
	device.InitDeployOptions()

	// вычислив адрес, нужно на него завести средства, чтобы вы
	walletAddress := device.CalcWalletAddress()

	// пополняем баланс wallet'a нового девайса
	giver := &everscale.Giver{
		Address: config.Get("giver.address"),
		Public:  config.Get("giver.public"),
		Secret:  config.Get("giver.secret"),
	}
	amount := 1_500_000_000
	log.Debugf("Giver: %s", giver.Address)
	log.Debug("Send Tokens from giver", "amount", amount, "from", giver.Address, "to", walletAddress, "amount", amount)
	err = giver.SendTokens("./contracts/giverv3.abi.json", walletAddress, amount)
	if err != nil {
		err = errors.Wrapf(err, "giver.SendTokens()")
		return
	}

	wait := 15 * time.Second
	log.Debugf("Wait %d seconds ...", wait.Seconds())
	time.Sleep(wait)

	// после всех сборок деплоим контракт
	log.Debug("Deploy ...")
	err = device.Deploy(data)
	if err != nil {
		err = errors.Wrapf(err, "device.Deploy(data)")
		return
	}

	// формируем ответ в формате json
	err = utils.JsonMapToStruct(data, &out)
	out.Address = dsm.EverAddress(walletAddress)

	log.Debug("device", out)
	return
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
