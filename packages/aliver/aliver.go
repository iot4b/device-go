package aliver

import (
	"time"

	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
)

var NodeHost string

func Run(s *coalago.Server, address string, aliveInterval time.Duration) {
	log.Info("run aliver")
	var retryErr int
	for {
		if NodeHost == "" {
			time.Sleep(time.Second)
			continue
		}
		aliveMessage := coalago.NewCoAPMessage(coalago.ACK, coalago.GET)
		aliveMessage.SetURIPath("/l")
		aliveMessage.SetURIQuery("a", address)
		_, err := s.Send(aliveMessage, NodeHost)
		if err != nil {
			log.Error(err, retryErr)
			retryErr++
			if retryErr > 10 {
				log.Error("retryErr > 10 - start registration")
				retryErr = 0
				if err := s.Refresh(); err != nil {
					log.Error("Refresh error:", err)
				} else {
					log.Info("Server refreshed")
					time.Sleep(2 * time.Second) // даём время новому listener'у запуститься
				}
			}
		} else {
			retryErr = 0
		}
		time.Sleep(aliveInterval)
	}
}
