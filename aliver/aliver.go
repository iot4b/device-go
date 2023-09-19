package aliver

import (
	"time"

	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
	"github.com/pkg/errors"
)

func Run(s *coalago.Server, address string, nodeHost string, aliveInterval time.Duration) {
	log.Info("run aliver")
	for {
		// если ставим после alive, то соединение с нодой не успевает инициироваться
		// todo пофиксить порядок запуска
		time.Sleep(aliveInterval)
		start := time.Now()
		err := alive(s, nodeHost, address)
		if err != nil {
			log.Error(err)
		}
		log.Debugf("time: %dµs node: %s address: %s", time.Since(start).Nanoseconds(), nodeHost, address)
	}
}

func alive(server *coalago.Server, nodeHost string, address string) error {
	aliveMessage := coalago.NewCoAPMessage(coalago.CON, coalago.GET)
	aliveMessage.SetURIPath("/a")
	aliveMessage.SetURIQuery("a", address)
	return errors.Wrap(server.SendToSocket(aliveMessage, nodeHost), "sendToSocket")
}
