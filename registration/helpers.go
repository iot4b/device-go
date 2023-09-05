package registration

import (
	"device-go/shared/config"
	log "device-go/shared/golog"
	"encoding/json"
	"github.com/coalalib/coalago"
	"math/rand"
	"time"
)

// ping - это получение info ноды,
// duration - время, затраченное на выполнение запроса
func ping(nodeHost string) (duration time.Duration, err error) {
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

// getEndpoints - получаем список нод с мастер ноды
func getEndpoints(masterNode string) (list []node, err error) {
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

// getMasterNode - получаем ноду из списка в конфигах
func getMasterNode() string {
	// получаем список нод по-умолчанию
	masterNodeList := config.List("masterNodes")

	// выбираем случайную ноду, чтобы одновременно на одну ноду не стучались все девайсы при инициализации, а было минимальное распределение
	randomIndex := rand.Intn(len(masterNodeList))
	masterNode := masterNodeList[randomIndex]
	// проверяем ноду на доступность, иначе пробуем следующую из списка
	_, err := ping(masterNode)
	if err != nil {
		log.Error(err)
		// удаляем ноду из списка и проходимся по оставшимся
		masterNodeList = append(masterNodeList[:randomIndex], masterNode[randomIndex+1:])
		for _, masterNode = range masterNodeList {
			_, err := ping(masterNode)
			if err != nil {
				log.Error(err)
			}
			if err == nil {
				break
			}
		}
	}
	return masterNode
}
