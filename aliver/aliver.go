package aliver

import (
	"time"

	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
)

/*

1 - запускаем элайвер
2 - если девайс прописался на ноде, то он шлет на нее элайв
4- если нода элайв реджектит, то мы зааускаем процесс регистрации
	если нрда не отвечает, то мы запускаем процесс регистрации
	если нода отвечает, и не реджектит, то все ок


*/

var NodeHost string

func Run(s *coalago.Server, address string, aliveInterval time.Duration) {
	log.Info("run aliver")
	var retryErr int = 0
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
				panic("reboot device") //переделать на нормальный рестарт сервиса
			}
		}
		time.Sleep(aliveInterval)
	}
}
