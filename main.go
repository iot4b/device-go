package main

import (
	"device-go/aliver"
	"device-go/cfg"
	"device-go/handlers"
	"device-go/models"

	"os"
	"os/signal"
	"syscall"

	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
)

func main() {
	handlers.Info = models.Info{
		Key:     cfg.GetString("publicKey"),
		Version: cfg.GetString("version"),
		Type:    cfg.GetString("type"),
		Vendor:  cfg.GetString("vendor"),
	}

	server := coalago.NewServer()
	server.GET("/info", handlers.GetInfo)
	server.POST("/cmd", handlers.ExecCmd)

	// начинаем слать alive пакеты, чтобы сохранять соединение для udp punching
	go aliver.Run(server, cfg.GetString("publicKey"), cfg.GetString("nodeHost"), cfg.GetTime("aliveInterval"))

	// стартуем сервер
	err := server.Listen(cfg.GetString("coapServerHost"))
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

func init() {
	cfg.Init("dev")
	log.Init(cfg.GetBool("debug"))
}
