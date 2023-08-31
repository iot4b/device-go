package registration

import (
	"device-go/helpers"
	"device-go/shared/config"
	"encoding/json"
	"errors"
	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
	"io"
	"os"
	"time"
)

type register struct {
	Key     string `json:"key"`
	Version string `json:"version"`
	Type    string `json:"type"`
	Vendor  string `json:"vendor"`
}

func Register(nodeHost, public, version, Type, vendor string) error {
	client := coalago.NewClient()

	payload := register{
		Key:     public,
		Version: version,
		Type:    Type,
		Vendor:  vendor,
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := coalago.NewCoAPMessage(coalago.CON, coalago.POST)
	msg.SetURIPath("/register")
	msg.SetStringPayload(string(bytes))
	resp, err := client.Send(msg, nodeHost)
	if err != nil {
		log.Error(err)
		return err
	}

	// сохраняем контракт девайса локально
	return helpers.SaveContractLocal(resp.Body)
}

func NodeList(masterNode string) (list []string, err error) {
	client := coalago.NewClient()

	msg := coalago.NewCoAPMessage(coalago.CON, coalago.GET)
	msg.SetURIPath("/endpoints")
	resp, err := client.Send(msg, masterNode)
	if err != nil {
		return nil, err
	}
	log.Debug(string(resp.Body))
	err = json.Unmarshal(resp.Body, &list)
	if err != nil {
		return nil, err
	}
	return
}

func Ping(nodeHost string) (duration time.Duration, err error) {
	start := time.Now()

	client := coalago.NewClient()

	msg := coalago.NewCoAPMessage(coalago.CON, coalago.GET)
	msg.SetURIPath("/info")
	resp, err := client.Send(msg, nodeHost)
	if err != nil {
		return
	}
	duration = time.Since(start)
	log.Debugf("node: %s, ping time: %d ms, %s", nodeHost, duration.Milliseconds(), string(resp.Body))
	return
}

func GetNode() (host string, needRegistration bool) {
	contractFile, err := os.Open(config.Get("device.contractFile"))
	// если файл не найден, то получаем ноду с минимальным пингом
	if err != nil && errors.Is(err, os.ErrNotExist) {

		// todo get rndm from master nodes
		nodeHost := config.Get("nodeHost")

		var list []string
		err := helpers.RoundRobin(func() error {
			var err error
			list, err = NodeList(nodeHost)
			return err
		}, 3*time.Second, 10)
		if err != nil {
			log.Fatal(err)
		}

		// check min ping time to host
		var lastTime time.Duration
		fasterHost := nodeHost
		for _, host := range list {
			t, err := Ping(host + config.Get("coapServerPort"))
			if err != nil {
				log.Error(err)
				continue
			}
			if lastTime > t || lastTime == 0 {
				lastTime = t
				fasterHost = host
			}
		}
		host = fasterHost
		needRegistration = true
		return
	}

	defer contractFile.Close()

	contract := map[string]interface{}{}
	// иначе читаем ноду из файла контракта
	data, err := io.ReadAll(contractFile)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &contract)
	if err != nil {
		log.Fatal(err)
	}
	// пингуем ноду из контракта
	host = contract["node"].(string)
	if _, err := Ping(host); err != nil {
		// если нода не пингуется, то удаляем текущий контракт и новую ноду выбираем для девайса
		err = os.Remove(config.Get("device.contractFile"))
		if err != nil {
			log.Error(err)
		}
		host, needRegistration = GetNode()
	}
	return
}
