package registration

import (
	"device-go/helpers"
	"device-go/shared/config"
	"encoding/json"
	"errors"
	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
	"io"
	"math/rand"
	"os"
	"time"
)

// Register - регистрируем устройство на ноде
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

// NodeList - получаем список нод с мастер ноды
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

// GetNode - определяем ноду к которой нужно подключиться по минимальному пингу.
// если есть файл контракта, то из него берет ноду.
// если девайс долго не был в сети, то может быть такая ситуация, что нода на которой он висел уже не активна,
// то повторяет процедуру регистрации
func GetNode() (host string, needRegistration bool) {
	contractFile, err := os.Open(config.Get("device.contractFile"))
	// если файл не найден, то получаем ноду с минимальным пингом
	if err != nil && errors.Is(err, os.ErrNotExist) {

		// get random from master nodes
		masterNodeList := config.List("masterNodes")

		randomIndex := rand.Intn(len(masterNodeList))
		masterNode := masterNodeList[randomIndex] + config.Get("coapServerPort")

		var list []string
		err := helpers.RoundRobin(func() error {
			var err error
			list, err = NodeList(masterNode)
			return err
		}, 3*time.Second, 10)
		if err != nil {
			log.Fatal(err)
		}

		// check min ping time to host
		var lastTime time.Duration
		fasterHost := masterNode
		for _, host := range list {
			t, err := ping(host + config.Get("coapServerPort"))
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
	if _, err := ping(host); err != nil {
		// если нода не пингуется, то удаляем текущий контракт и новую ноду выбираем для девайса
		err = os.Remove(config.Get("device.contractFile"))
		if err != nil {
			log.Error(err)
		}
		host, needRegistration = GetNode()
	}
	return
}
