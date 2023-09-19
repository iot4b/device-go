package main

import (
	"device-go/aliver"
	"device-go/cryptoKeys"
	"device-go/dsm"
	"device-go/handlers"
	"device-go/registration"
	"device-go/shared"
	"device-go/shared/config"
	"device-go/storage"
	"fmt"
	"github.com/coalalib/coalago"
	"github.com/jinzhu/copier"
	log "github.com/ndmsystems/golog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//TODO  при старте девайса надо скачать смартконтракт девайса и сохранить локально.
// В смарт контракте прописаны ключи которые имеют право присылать команды,  смарт вендора, из которого берем имя вендора для конфига
// по клюбчам проверяем что команда подписана тем ключем, который стоит в разрешенных, и тогда выполняем ее.

func main() {
	// инициируем ключи девайса. если есть файл, то читаем из него, если нет, то генерим новый
	// для ключей используется алгоритм ed25519
	cryptoKeys.Init()
	storage.Init(
		config.Get("localFiles.contractData"),
		config.Get("vendorName"),
		config.Get("type"),
		config.Get("version"),
		"")

	var nodeHost string // nodeHost нужен, чтобы передать его в alive
	var registeredDevice *dsm.DeviceContract
	var err error

	for {
		// регистрируем устройство на ноде. в ответ приходит нода, к которой получилось подключиться
		// если ошибка, то повторяем цикл регистрации
		registeredDevice, nodeHost, err = registration.Register(
			// получаем список нод по-умолчанию
			config.List("masterNodes"),
			storage.Get().Address,
			storage.Get().VendorAddress,

			cryptoKeys.KeyPair.PublicStr(),
			storage.Get().Version,
			storage.Get().Type,
			storage.Get().VendorData)

		if err == nil {
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

	// копируем данные девайса в Info
	copier.Copy(&shared.Info, storage.Get())
	shared.Info.Key = cryptoKeys.KeyPair.PublicStr()
	shared.Info.Uptime = time.Since(shared.Info.RunFrom).String()

	log.Debug("device info: %+v", shared.Info)

	// сервер для запросов от клиентов и нод
	server := coalago.NewServer()
	server.GET("/info", handlers.GetInfo)
	server.POST("/cmd", handlers.ExecCmd)

	// начинаем слать alive пакеты, чтобы сохранять соединение для udp punching
	go aliver.Run(server, storage.Get().Address.String(), nodeHost, config.Time("aliveInterval"))
	// стартуем сервер
	err = server.Listen(":" + config.Get("device.port"))
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

// инитим конфиги и logger
func init() {
	if len(os.Args) < 1 {
		fmt.Println(`Usage: server [env]`)
		fmt.Println("Not enough arguments. Use defaults : dev")
		os.Exit(0)
	}
	config.Init(os.Args[1])
	log.Init(config.Bool("debug"))
}
