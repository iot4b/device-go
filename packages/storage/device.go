package storage

import (
	"bufio"
	"device-go/packages/config"
	"device-go/packages/dsm"
	"device-go/packages/utils"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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

	NodeIpPort string `json:"nodeIpPort,omitempty"`
}

var (
	Device   device
	filePath string
)

func Init(path, initFile, elector, vendor, deviceAPI, dType, version string) {
	filePath = filepath.Join(utils.GetFilesDir(), path)

	log.Debug(path, elector, vendor, deviceAPI, dType, version)

	var err error
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
			name, group, owner := promptUserData()
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

func WaitForData() {
	filePath = filepath.Join(utils.GetFilesDir(), config.Get("localFiles.contract"))
	if !utils.FileExists(filePath) {
		log.Info("Waiting for contract data file...")
		for {
			time.Sleep(time.Second)
			if utils.FileExists(filePath) {
				break
			}
		}
	}
	Update()
}

// Update Device data from file
func Update() (err error) {
	Device, err = read(filePath)
	return
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

func promptUserData() (name, group, owner string) {
	reader := bufio.NewReader(os.Stdin)
	name, _ = os.Hostname()
	fmt.Printf("Enter device name [%s]\n", name)
	for {
		name, _ = reader.ReadString('\n')
		name = strings.TrimSpace(name)
		if name != "" {
			break
		}
		fmt.Println("Device name is required, please try again")
	}
	fmt.Println("Enter device group address")
	for {
		group, _ = reader.ReadString('\n')
		group = strings.TrimSpace(group)
		if utils.MatchRegex(`^0:[0-9a-fA-F]{64}$`, group) {
			break
		}
		fmt.Println("Please enter group address in a format:")
		fmt.Println("0:0000000000000000000000000000000000000000000000000000000000000000")
	}
	fmt.Println("Enter owner public key")
	for {
		owner, _ = reader.ReadString('\n')
		owner = strings.TrimSpace(owner)
		if utils.MatchRegex(`^0x[0-9a-fA-F]{64}$`, owner) {
			break
		}
		fmt.Println("Please enter public key in a format:")
		fmt.Println("0x0000000000000000000000000000000000000000000000000000000000000000")
	}
	return
}
