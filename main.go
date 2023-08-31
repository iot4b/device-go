package main

import (
	"device-go/aliver"
	"device-go/crypto"
	"device-go/handlers"
	"device-go/helpers"
	"device-go/models"
	"device-go/registration"
	"device-go/shared/config"
	"fmt"
	"time"

	"os"
	"os/signal"
	"syscall"

	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
)

//TODO  при старте девайса надо скачать смартконтракт девайса и сохранить локально.
// В смарт контракте прописаны ключи которые имеют право присылать команды,  смарт вендора, из которого берем имя вендора для конфига
// по клюбчам проверяем что команда подписана тем ключем, который стоит в разрешенных, и тогда выполняем ее.

func main() {
	crypto.Init()

	// todo get rndm from master nodes
	nodeHost := config.Get("nodeHost")

	var list []string
	err := helpers.RoundRobin(func() error {
		var err error
		list, err = registration.NodeList(nodeHost)
		return err
	}, 3*time.Second, 10)
	if err != nil {
		log.Fatal(err)
	}

	// check min ping time to host
	var lastTime time.Duration
	fasterHost := nodeHost
	for _, host := range list {
		t, err := registration.Ping(host + config.Get("coapServerPort"))
		if err != nil {
			log.Error(err)
			continue
		}
		if lastTime > t || lastTime == 0 {
			lastTime = t
			fasterHost = host
		}
	}

	err = helpers.RoundRobin(func() error {
		return registration.Register(fasterHost + config.Get("coapServerPort"))
	}, 3*time.Second, 10)
	if err != nil {
		log.Fatal(err)
	}

	handlers.Info = models.Info{
		Key:     crypto.KeyPair.PublicStr(),
		Version: config.Get("version"),
		Type:    config.Get("type"),
		Vendor:  config.Get("vendor"),
	}

	//server := coalago.NewServerWithPrivateKey([]byte(crypto.KeyPair.Secret))
	server := coalago.NewServer()
	server.GET("/info", handlers.GetInfo)
	server.POST("/cmd", handlers.ExecCmd)

	// начинаем слать alive пакеты, чтобы сохранять соединение для udp punching
	go aliver.Run(server, crypto.KeyPair.PublicStr(), config.Get("nodeHost"), config.Time("aliveInterval"))

	// стартуем сервер
	err = server.Listen(config.Get("coapServerHost"))
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
