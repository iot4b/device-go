package registration

import (
	"encoding/json"
	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
	"github.com/pkg/errors"
	"time"
)

type node struct {
	IpPort  string `json:"ipPort"`
	Account string `json:"account"`
}

// Register - регистрируем устройство на ноде. Возвращает адрес ноды
func Register(public, version, Type, vendor string) (string, error) {
	masterNode := getMasterNode()

	// определив мастер ноду, получаем с нее список нод
	list, err := getEndpoints(masterNode)
	if err != nil {
		return "", errors.Wrap(err, "getEndpoints")
	}

	// check min ping time to host
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

	client := coalago.NewClient()

	payload := register{
		Key:     public,
		Version: version,
		Type:    Type,
		Vendor:  vendor,
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		return "", errors.Wrap(err, "json.Marshal(payload)")
	}

	msg := coalago.NewCoAPMessage(coalago.CON, coalago.POST)
	msg.SetURIPath("/register")
	msg.SetStringPayload(string(bytes))
	_, err = client.Send(msg, fasterHost)
	if err != nil {
		return "", errors.Wrap(err, "client.Send")
	}
	return fasterHost, nil
}
