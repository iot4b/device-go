package api

import (
	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
)

func Start() {
	coalaServer := coalago.NewServer()

	coalaServer.GET("/info", getInfo)
	coalaServer.POST("/cmd", execCmd)

	log.Fatal(coalaServer.Listen("127.0.0.1:5683"))
}
