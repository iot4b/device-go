package registration

import (
	"device-go/crypto"
	"device-go/dsm"
	"device-go/everscale"
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
func Register(masterNodes []string, address, vendorAddress dsm.EverAddress, public, version, Type, vendorData string) (*dsm.DeviceContract, string, error) {
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

	// формируем запрос на регистрацию
	client := coalago.NewClient()
	msg := coalago.NewCoAPMessage(coalago.CON, coalago.POST)
	msg.SetURIPath("/register")
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

	// set node to blockchain
	input := map[string]interface{}{}
	input["value"] = fasterAddress
	s := everscale.NewSigner(crypto.Keys.PublicSign, crypto.Keys.Secret)
	r, err := everscale.Execute("Device", string(address), "setNode", input, s)
	if err != nil {
		log.Debug(err, "setNode: "+string(r), input)
		return nil, "", errors.Wrap(err, "setNode")
	}
	result.Node = dsm.EverAddress(fasterAddress)

	log.Debugw("Register result", "RegisteredDevice", result, "fasterHost", fasterHost)

	return &result, fasterHost, nil
}
