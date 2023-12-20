package aliver

import (
	"time"

	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
	"github.com/pkg/errors"
)

/*

1 - запускаем элайвер
2 - если девайс прописался на ноде, то он шлет на нее элайв
4- если нода элайв реджектит, то мы зааускаем процесс регистрации
	если нрда не отвечает, то мы запускаем процесс регистрации
	если нода отвечает, и не реджектит, то все ок


*/

func Run(s *coalago.Server, address string, nodeHost string, aliveInterval time.Duration) {
	log.Info("run aliver")
	for {
		// если ставим после alive, то соединение с нодой не успевает инициироваться
		time.Sleep(aliveInterval)
		//start := time.Now()
		err := alive(s, nodeHost, address)
		if err != nil {
			log.Error(err)
			//todo если не удалось отправить сообщение, то pзапускаем процесс переподключения, если накопилось 10 ошибок

		}
		//log.Infof("time: %dµs node: %s address: %s", time.Since(start).Nanoseconds(), nodeHost, address)
	}
}

func alive(server *coalago.Server, nodeHost string, address string) error {
	aliveMessage := coalago.NewCoAPMessage(coalago.ACK, coalago.GET)
	aliveMessage.SetURIPath("/l")
	aliveMessage.SetURIQuery("a", address)
	//todo  - переписать на нормальную отправку
	return errors.Wrap(server.SendToSocket(aliveMessage, nodeHost), "sendToSocket")
}
