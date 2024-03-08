package registration

import (
	"device-go/aliver"
	"device-go/crypto"
	"device-go/dsm"
	"device-go/everscale"
	"device-go/shared/config"
	"device-go/storage"
	"encoding/json"
	"github.com/coalalib/coalago"
	"github.com/jinzhu/copier"
	log "github.com/ndmsystems/golog"
	"github.com/pkg/errors"
	"time"
)

type registerDeviceResp struct {
	Address dsm.EverAddress `json:"a,omitempty"` //ever SC address текущего Device
	Node    dsm.EverAddress `json:"n,omitempty"` //ever SC address Node, с которой девайс создал последнее соединение
	Elector dsm.EverAddress `json:"e,omitempty"` //ever SC адрес Elector'a, который обслуживает сеть нод для текущего девайса
	Vendor  dsm.EverAddress `json:"v,omitempty"` //ever SC address производителя текущего девайса

	Stat bool `json:"s,omitempty"` // нужно ли девайсу слать статистику
}

// Register - регистрируем устройство на ноде.
// Возвращает ip:port ноды
func Register() (*dsm.DeviceContract, string, error) {
	log.Debug("Register")

	masterNodes := config.List("masterNodes")
	address := storage.Get().Address
	vendorAddress := storage.Get().Vendor

	public := crypto.Keys.PublicSign
	version := storage.Get().Version
	Type := storage.Get().Type
	vendorData := storage.Get().VendorData

	// получаем список доступных нод с рандомной мастер ноды
	masterNode, list, err := endpointList(masterNodes)
	if err != nil {
		return nil, "", errors.Wrap(err, "getEndpoints")
	}
	log.Debug("endpoints: %+v", list)

	// перебираем ноды и определяем самый низкий ping, далее используем эту ноду для регистрации и поддержания соединения
	var lastTime time.Duration
	fasterHost := masterNode
	fasterAddress := ""
	log.Debug("fasterHost before ping: " + fasterHost)
	for _, host := range list {
		t, err := ping(host.IpPort)
		if err != nil {
			if host.IpPort != "240.0.0.0:65535" { // non-existent node
				log.Error(err)
			}
			continue
		}
		if lastTime > t || lastTime == 0 {
			lastTime = t
			fasterHost = host.IpPort
			fasterAddress = host.Account
		}
	}
	log.Info("Registering device on node:", fasterHost)

	payload, err := json.Marshal(registerRequest{
		Address:    address,
		Vendor:     vendorAddress,
		Key:        public,
		Version:    version,
		Type:       Type,
		VendorData: vendorData,
	})
	if err != nil {
		return nil, "", errors.Wrap(err, "json.Marshal(payload)")
	}
	log.Debug("registerRequest: " + string(payload) + " address: " + address.String())

	// update device node if differs
	if storage.Get().Node != dsm.EverAddress(fasterAddress) {
		err = everscale.Device.SetNode(dsm.EverAddress(fasterAddress))
		if err != nil {
			log.Error("everscale.Device.SetNode:", err)
			return nil, "", err
		}
	}

	// формируем запрос на регистрацию
	client := coalago.NewClient()
	msg := coalago.NewCoAPMessage(coalago.CON, coalago.POST)
	msg.SetURIPath("/register")
	msg.Timeout = config.Time("timeout.coala")
	// если девайс знает свой адрес контракта, то передаем ...?a= ...
	if len(address) > 0 {
		msg.SetURIQuery("a", string(address))
	}
	msg.SetStringPayload(string(payload))

	// отправляем запрос на регистрацию
	resp, err := client.Send(msg, fasterHost)
	if err != nil {
		return nil, "", errors.Wrap(err, "client.Send")
	}

	// парсим ответ и обновляем локальный дамп контракта
	registerResp := registerDeviceResp{}
	err = json.Unmarshal(resp.Body, &registerResp)
	if err != nil {
		log.Debug("registerResponse: "+string(resp.Body), "code: "+string(resp.Code))
		return nil, "", errors.Wrap(err, "json.Unmarshal(resp.Body, &registerResp)")
	}
	log.Debug("registerResponse: " + string(resp.Body))

	// копируем актуальные поля
	result := dsm.DeviceContract{}
	copier.Copy(&result, registerResp)

	result.Node = dsm.EverAddress(fasterAddress)

	log.Debugw("Register result", "RegisteredDevice", result, "fasterHost", fasterHost)

	return &result, fasterHost, nil
}

// Repeat registration every 15m
func Repeat() {
	log.Info("Repeat registration")
	for {
		// update device data from smartcontract
		device, err := everscale.Device.Get()
		if err != nil {
			log.Error("everscale.Device.Get:", err)
			time.Sleep(3 * time.Second)
			continue
		}
		storage.Set(device)
		storage.WriteToLocalStorage(config.Get("localFiles.contract"), device)

		_, nodeHost, err := Register()
		if err != nil {
			log.Error(err)
			time.Sleep(3 * time.Second)
			continue
		}
		aliver.NodeHost = nodeHost
		time.Sleep(15 * time.Minute)
	}
}
