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
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
)

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

	var nodeHost string // nodeHost нужен, чтобы передать его в alive
	var registeredDevice *dsm.DeviceContract
	var err error

	for {
		// регистрируем устройство на ноде. в ответ приходит нода, к которой получилось подключиться
		// если ошибка, то повторяем цикл регистрации
		registeredDevice, nodeHost, err = registration.Register()
		if err == nil {
			aliver.NodeHost = nodeHost
			break
		}
		log.Error(err)
		time.Sleep(config.Time("timeout.registerRepeat"))
	}

	// если регистрация прошла успешно, то нужно обновить данные о текущем девайсе в локальном хранилище
	err = storage.Update(registeredDevice)
	if err != nil {
		log.Fatal(err)
	}

	// add owner if passed via -addOwner flag
	var owner string
	flag.StringVar(&owner, "addOwner", "", "add owner public key to device contract")
	if len(owner) > 0 {
		log.Info("addOwner:", owner)
		if storage.IsOwner(owner) {
			log.Warning("already an owner, skipping")
		} else if err := everscale.Device.AddOwner(owner); err != nil {
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

	// начинаем слать alive пакеты, чтобы сохранять соединение для udp punching
	go aliver.Run(server, storage.Get().Address.String(), config.Time("timeout.alive"))

	// стартуем сервер
	port := config.Get("port.device")
	if len(os.Args) > 2 {
		port = os.Args[2]
	}
	err = server.Listen(":" + port)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

// инитим конфиги и logger
func init() {
	if len(os.Args) < 2 {
		fmt.Println(`Usage: server [env]`)
		fmt.Println("Not enough arguments. Use defaults : dev")
		os.Exit(0)
	}
	config.Init(os.Args[1])
	log.Init(config.Bool("debug"))
}
