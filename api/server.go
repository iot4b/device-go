package api

import (
	"device-go/models"
	"time"

	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
	"github.com/pkg/errors"
)

var (
	info models.Info
)

func NewServer(i models.Info) *coalago.Server {
	info = i
	server := coalago.NewServer()

	server.GET("/info", GetInfo)
	server.POST("/cmd", ExecCmd)

	return server
}

func Start(s *coalago.Server, host string) {
	log.Info("api start")
	log.Fatal(s.Listen(host))
}

func RunAlive(publicKey, nodeHost string, aliveInterval time.Duration, s *coalago.Server) {
	log.Info("run alive")
	//c := client.New(nodeHost, []byte(publicKey))
	//отсылаем alive на ноду
	for {
		// если ставим после alive, то соединение с нодой не успевает инициироваться
		// todo пофиксить порядок запуска
		time.Sleep(aliveInterval)
		//c.SendAlive()
		if err := alive(s, nodeHost, publicKey); err != nil {
			log.Error(err)
		}
	}
}

func alive(server *coalago.Server, nodeHost, publicKey string) error {
	aliveMessage := coalago.NewCoAPMessage(coalago.CON, coalago.GET)
	aliveMessage.SetURIPath("/live")
	aliveMessage.SetURIQuery("key", publicKey)

	log.Info("send alive", aliveMessage.Payload.String(), nodeHost, publicKey)
	return errors.Wrap(server.SendToSocket(aliveMessage, nodeHost), "sendToSocket")
}
