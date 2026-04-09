package registration

import (
	"device-go/packages/config"
	"device-go/packages/crypto"
	"device-go/packages/dsm"
	"device-go/packages/storage"
	"encoding/json"
	"fmt"

	"github.com/coalalib/coalago"
	"github.com/jinzhu/copier"
	log "github.com/ndmsystems/golog"
)

// Register - регистрируем устройство на ноде.
// Возвращает ip:port ноды
func Register() error {
	log.Debug("Register")

	fasterHost, fasterAddress, err := selectFastestNode()
	if err != nil {
		return err
	}
	log.Info("Registering device on node:", fasterHost)

	var req registerRequest
	copier.Copy(&req, storage.Device)
	req.PublicSign = crypto.Keys.PublicSign
	req.PublicNacl = crypto.Keys.PublicNacl
	req.Events = true

	payload, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}
	log.Debugf("registerRequest: %s address: %s", payload, storage.Device.Address)

	// формируем запрос на регистрацию
	client := coalago.NewClient()
	msg := coalago.NewCoAPMessage(coalago.CON, coalago.POST)
	msg.SetURIPath("/register")
	msg.Timeout = config.Time("timeout.coala")
	// если девайс знает свой адрес контракта, то передаем ...?a= ...
	if len(storage.Device.Address) > 0 {
		msg.SetURIQuery("a", string(storage.Device.Address))
	}
	msg.SetStringPayload(string(payload))

	// отправляем запрос на регистрацию
	resp, err := client.Send(msg, fasterHost)
	if err != nil {
		return fmt.Errorf("client.Send: %w", err)
	}
	if resp.Code.IsCommonError() || resp.Code.IsInternalError() {
		return fmt.Errorf("client.Send: %s (%v)", resp.Body, resp.Code)
	}

	// парсим ответ и обновляем локальный дамп контракта
	registerResp := registerResponse{}
	err = json.Unmarshal(resp.Body, &registerResp)
	if err != nil {
		return fmt.Errorf("json.Unmarshal(%s): %v", resp.Body, err)
	}
	log.Debug("registerResponse: " + string(resp.Body))

	// копируем актуальные поля
	copier.Copy(&storage.Device, registerResp)
	if storage.Device.Node == "" && fasterAddress != "" {
		storage.Device.Node = dsm.EverAddress(fasterAddress)
	}
	storage.Device.NodeIpPort = fasterHost

	// update local data
	storage.Save()

	log.Infow("Register result", "RegisteredDevice", storage.Device, "nodeHost", fasterHost)

	return nil
}
