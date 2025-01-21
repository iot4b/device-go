package registration

import (
	"device-go/packages/config"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
)

// endpointList - получаем список доступных нод для подключения
// смарт-контракт электора хранит список доступных нод, с которыми можно установить соединение. при запросе на одну из
// нод, можем получить полный список всех активных нод в системе
//
// выбранная нода в дальнейшем будет обслуживать девайс и устанавливать соединения с Client->Device
//
// в конфигах есть список мастер-нод, которые по-умолчанию задаются вместе с прошивкой девайса
// они постоянные и не меняются
//
// todo сделать привязку к доменам, чтобы мастер ноды, вне зависимости от их ip, были всегда доступны для запроса
// из masterNodeList выбирает случайным образом одну из мастер нод
func endpointList(masterNodeList []string) (masterNode string, list []node, err error) {
	// выбираем случайную ноду, чтобы одновременно на одну ноду не стучались все девайсы при инициализации,
	// а было минимальное распределение
	// перемешиваем список нод для случайного выбора
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(masterNodeList), func(i, j int) {
		masterNodeList[i], masterNodeList[j] = masterNodeList[j], masterNodeList[i]
	})

	// пробегаемся по нодам и выбираем первую доступную
	for _, masterNode = range masterNodeList {
		// получаем с нее список нод
		list, err = getEndpoints(masterNode)
		if err == nil {
			log.Debugf("masterNode: %s", masterNode)
			return
		}
		log.Errorf("masterNode: %s, err: %s", masterNode, err.Error())
	}
	return
}

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

// getEndpoints - получаем список доступных нод с указанной ноды
func getEndpoints(node string) (list []node, err error) {
	client := coalago.NewClient()

	msg := coalago.NewCoAPMessage(coalago.CON, coalago.GET)
	msg.SetURIPath("/endpoints")
	msg.Timeout = config.Time("timeout.coala")

	resp, err := client.Send(msg, node)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp.Body, &list)
	if err != nil {
		return nil, err
	}
	return
}
