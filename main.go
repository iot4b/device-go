package main

import (
	"device-go/aliver"
	"device-go/handlers"
	"device-go/models"
	"device-go/shared/config"
	"fmt"

	"os"
	"os/signal"
	"syscall"

	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
)

func main() {

	//TODO  при старте девайса надо скачать смартконтракт девайса и сохранить локально.
	// В смарт контракте прописаны ключи которые имеют право присылать команды,  смарт вендора, из которого берем имя вендора для конфига
	// по клюбчам проверяем что команда подписана тем ключем, который стоит в разрешенных, и тогда выполняем ее.
	//

	handlers.Info = models.Info{
		Key:     config.Get("publicKey"),
		Version: config.Get("version"),
		Type:    config.Get("type"),
		Vendor:  config.Get("vendor"),
	}

	server := coalago.NewServer()
	server.GET("/info", handlers.GetInfo)
	server.POST("/cmd", handlers.ExecCmd)

	// начинаем слать alive пакеты, чтобы сохранять соединение для udp punching
	go aliver.Run(server, config.Get("publicKey"), config.Get("nodeHost"), config.Time("aliveInterval"))

	// стартуем сервер
	err := server.Listen(config.Get("coapServerHost"))
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
