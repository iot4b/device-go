package main

import (
	"device-go/aliver"
	"device-go/crypto"
	"device-go/dsm"
	"device-go/everscale"
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
	everscale.Init(config.List("everscale.endpoints"))
	defer everscale.Destroy()

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
		config.List("everscale.owners"))

	everscale.Device.Address = storage.Get().Address

	// add new owner if passed via -addOwner flag
	if len(newOwner) > 0 {
		log.Info("add owner:", newOwner)
		if storage.IsOwner(newOwner) {
			log.Warning("already an owner, skipping")
		} else if err := everscale.Device.AddOwner(newOwner); err != nil {
			log.Fatal(err)
		}
		log.Info("new owner added")
	}

	// сервер для запросов от клиентов и нод
	server := coalago.NewServer()
	server.GET("/info", handlers.GetInfo)
	server.POST("/cmd", handlers.ExecCmd)
	server.GET("/confirm", handlers.Confirm)
	server.POST("/update", handlers.Update)

	var nodeHost string // nodeHost нужен, чтобы передать его в alive
	var registeredDevice *dsm.DeviceContract
	var err error

	for {
		// регистрируем устройство на ноде. в ответ приходит нода, к которой получилось подключиться
		// если ошибка, то повторяем цикл регистрации
		registeredDevice, nodeHost, err = registration.Register()
		if err == nil {
			// если регистрация прошла успешно, то нужно обновить данные о текущем девайсе в локальном хранилище
			err = storage.Update(registeredDevice)
			if err != nil {
				log.Fatal(err)
			}

			// стартуем сервер
			go listen(server)

			// начинаем слать alive пакеты, чтобы сохранять соединение для udp punching
			aliver.NodeHost = nodeHost
			go aliver.Run(server, storage.Get().Address.String(), config.Time("timeout.alive"))

			time.Sleep(15 * time.Minute)
			go registration.Repeat()

			break
		}
		log.Error(err)
		time.Sleep(config.Time("timeout.registerRepeat"))
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
