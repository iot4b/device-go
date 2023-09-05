package aliver

import (
	"time"

	log "device-go/shared/golog"
	"github.com/coalalib/coalago"
	"github.com/pkg/errors"
)

func Run(s *coalago.Server, publicKey, nodeHost string, aliveInterval time.Duration) {
	log.Info("run alive")
	for {
		// если ставим после alive, то соединение с нодой не успевает инициироваться
		// todo пофиксить порядок запуска
		time.Sleep(aliveInterval)
		err := alive(s, nodeHost, publicKey)
		if err != nil {
			log.Error(err)
		}
	}
}

func alive(server *coalago.Server, nodeHost, publicKey string) error {
	aliveMessage := coalago.NewCoAPMessage(coalago.CON, coalago.GET)
	aliveMessage.SetURIPath("/live")
	aliveMessage.SetURIQuery("key", publicKey)
	log.Debug(aliveMessage.Payload.String(), nodeHost, publicKey)
	return errors.Wrap(server.SendToSocket(aliveMessage, nodeHost), "sendToSocket")
}
