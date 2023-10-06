package registration

import (
	"device-go/storage"
	"encoding/json"
	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
	"github.com/pkg/errors"
	"time"
)

// Register - регистрируем устройство на ноде. Возвращает адрес ноды
func Register(masterNodes []string, public, version, Type, vendor string) (string, error) {
	// получаем список доступных нод с рандомной мастер ноды
	masterNode, list, err := endpointList(masterNodes)
	if err != nil {
		return "", errors.Wrap(err, "getEndpoints")
	}
	log.Debug("endpoints: %+v", list)

	// перебираем ноды и определяем самый низкий ping, далее используем эту ноду для регистрации и поддержания соединения
	var lastTime time.Duration
	fasterHost := masterNode
	for _, host := range list {
		t, err := ping(host.IpPort)
		if err != nil {
			log.Error(err)
			continue
		}
		if lastTime > t || lastTime == 0 {
			lastTime = t
			fasterHost = host.IpPort
		}
	}

	payload, err := json.Marshal(register{
		Address: storage.Contract.Address,
		Version: version,
		Type:    Type,
		Vendor:  vendor,
	})
	if err != nil {
		return "", errors.Wrap(err, "json.Marshal(payload)")
	}

	// отправляем запрос на регистрацию
	client := coalago.NewClient()
	msg := coalago.NewCoAPMessage(coalago.CON, coalago.POST)
	msg.SetURIPath("/register")
	msg.SetStringPayload(string(payload))

	res, err := client.Send(msg, fasterHost)
	if err != nil {
		return "", errors.Wrap(err, "client.Send")
	}

	log.Debug(string(res.Body))

	if len(storage.Contract.Address) == 0 {
		if err = json.Unmarshal(res.Body, &storage.Contract); err != nil {
			return "", errors.Wrap(err, "json.Unmarshal")
		}
		storage.Contract.Save()
	}

	return fasterHost, nil
}
