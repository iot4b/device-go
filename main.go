package main

import (
	"device-go/api"
	"device-go/cfg"
	"device-go/models"
	log "github.com/ndmsystems/golog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	info := models.Info{
		Key:     cfg.GetString("publicKey"),
		Version: cfg.GetString("version"),
		Type:    cfg.GetString("type"),
		Vendor:  cfg.GetString("vendor"),
	}

	// инициируем CoaP сервер
	coapServer := api.NewServer(info)

	// стартуем сервер
	go api.Start(coapServer, cfg.GetString("coapServerHost"))

	// начинаем слать alive пакеты, чтобы сохранять соединение для udp punching
	go api.RunAlive(cfg.GetString("publicKey"),
		cfg.GetString("nodeHost"),
		cfg.GetTime("aliveInterval"),
		coapServer)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

func init() {
	log.Init("en", "dev", "test")
	cfg.Init("dev")
}
