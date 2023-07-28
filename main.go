package main

import (
	"device-go/api"
	"device-go/cfg"
	"device-go/client"
	"device-go/models"
	"device-go/workers"
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
	go api.Start(info, []byte(cfg.GetString("privateKey")))

	cl := client.New(cfg.GetString("node"), []byte(cfg.GetString("publicKey")))
	go workers.RunAlive(cl)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

func init() {
	log.Init("en", "dev", "test")
	cfg.Init("dev")
}
