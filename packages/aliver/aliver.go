package aliver

import (
	"time"

	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
)

var NodeHost string

func Run(s *coalago.Server, address string, aliveInterval time.Duration) {
	log.Info("run aliver")
	var retryErr = 0
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
				//todo запуск процесса регистрации
				retryErr = 0

				//restart service
				//TODO это костыль надо чтобы коала помнила адрес и перезапусклась на нем после рестарта
				//flag.StringVar(&port, "port", port, "override default coala port")
				//flag.Parse()
				//err = s.Listen(":" + port)
				//if err != nil {
				//	log.Panic(err)
				//}
			}
		}
		time.Sleep(aliveInterval)
	}
}
