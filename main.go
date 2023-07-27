package main

import (
	"device-go/api"
	"device-go/client"
	"device-go/workers"
	log "github.com/ndmsystems/golog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Init("en", "dev", "test")

	go api.Start()

	cl := client.New()
	go workers.RunAlive(cl)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
