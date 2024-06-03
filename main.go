package main

import (
	"device-go/aliver"
	"device-go/crypto"
	"device-go/events"
	"device-go/handlers"
	"device-go/registration"
	"device-go/shared/config"
	"device-go/storage"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
)

var port, newOwner string

//TODO  при старте девайса надо скачать смартконтракт девайса и сохранить локально.
// В смарт контракте прописаны ключи которые имеют право присылать команды,  смарт вендора, из которого берем имя вендора для конфига
// по клюбчам проверяем что команда подписана тем ключем, который стоит в разрешенных, и тогда выполняем ее.

func main() {
	// инициируем ключи девайса. если есть файл, то читаем из него, если нет, то генерим новый
	// для ключей используется алгоритм ed25519
	crypto.Init(config.Get("localFiles.keys"))
	storage.Init(
		config.Get("localFiles.contract"),
		config.Get("everscale.elector"),
		config.Get("everscale.vendor.address"),
		config.Get("everscale.vendor.name"),
		config.Get("everscale.vendor.data"),
		config.Get("info.type"),
		config.Get("info.version"),
		config.Map("everscale.owners"),
		config.Get("everscale.group"))

	// сервер для запросов от клиентов и нод
	server := coalago.NewServer()
	server.GET("/info", handlers.GetInfo)
	server.POST("/cmd", handlers.ExecCmd)
	server.POST("/update", handlers.Update)
	server.POST("/sign", handlers.Sign)

	for {
		// регистрируем устройство на ноде
		// если ошибка, то повторяем цикл регистрации
		err := registration.Register()
		if err == nil {
			// стартуем сервер
			go listen(server)

			// delay for alive correct work
			time.Sleep(time.Second)

			// начинаем слать alive пакеты, чтобы сохранять соединение для udp punching
			go aliver.Run(server, storage.Device.Address.String(), config.Time("timeout.alive"))

			if storage.Device.Events {
				// delay after first alive to store ip:port and send event
				time.Sleep(time.Second)
				events.Send(new(events.Register))
			}

			time.Sleep(config.Time("timeout.registerRepeat"))
			go registration.Repeat()

			break
		}
		log.Error(err)
		time.Sleep(3 * time.Second)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

// инитим конфиги и logger
func init() {
	var env string
	flag.StringVar(&env, "env", "dev", "set environment")
	flag.StringVar(&port, "port", "5683", "set coala port")
	flag.StringVar(&newOwner, "addOwner", "", "add new owner public key to device contract")
	flag.Parse()

	config.Init(env)
	log.Init(config.Bool("debug"))
}

func listen(server *coalago.Server) {
	if err := server.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
