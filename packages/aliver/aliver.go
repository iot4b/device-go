package aliver

import (
	"device-go/packages/storage"
	"time"

	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
)

func Run(s *coalago.Server, aliveInterval time.Duration) {
	log.Info("run aliver")
	var retryErr int
	for {
		aliveMessage := coalago.NewCoAPMessage(coalago.ACK, coalago.GET)
		aliveMessage.SetSchemeCOAPS()
		aliveMessage.SetURIPath("/l")
		aliveMessage.SetURIQuery("a", storage.Device.Address.String())
		_, err := s.Send(aliveMessage, storage.Device.NodeIpPort)
		if err != nil {
			log.Error(err, retryErr)
			retryErr++
			if retryErr > 10 {
				log.Error("retryErr > 10 - start registration")
				return
			}
		} else {
			retryErr = 0
		}
		time.Sleep(aliveInterval)
	}
}
