package registration

import (
	"device-go/shared/config"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
)

// Register - регистрируем устройство на ноде. Возвращает адрес ноды
func Register(public, version, Type, vendor string) (string, error) { //TODO возвращать ошибку.   И гонять регистер по круту с интервалом timeout.registerRepeat
	// если файл не найден, то получаем ноду с минимальным пингом

	// get random from master nodes
	masterNodeList := config.List("masterNodes") //TODO а если матер не ответил?  переделать на цикл, пока хоть ктото не ответит

	randomIndex := rand.Intn(len(masterNodeList))
	masterNode := masterNodeList[randomIndex]

	var list []string
	list, err := getEndpoints(masterNode)
	if err != nil {
		log.Fatal(err)
	}

	// check min ping time to host
	var lastTime time.Duration
	fasterHost := masterNode
	for _, host := range list {
		t, err := ping(host + config.Get("coapServerPort")) //TODO порт прийдет с айпишником
		if err != nil {
			log.Error(err)
			continue
		}
		if lastTime > t || lastTime == 0 {
			lastTime = t
			fasterHost = host
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
		log.Error(err)
		return "", err
	}

	msg := coalago.NewCoAPMessage(coalago.CON, coalago.POST)
	msg.SetURIPath("/register")
	msg.SetStringPayload(string(bytes))
	_, err = client.Send(msg, fasterHost)
	if err != nil {
		log.Error(err)
		return "", err
	}
	return fasterHost, nil
}
