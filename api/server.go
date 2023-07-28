package api

import (
	"device-go/models"
	"github.com/coalalib/coalago/coalaServer"
	log "github.com/ndmsystems/golog"
)

var (
	info       models.Info
	privateKey []byte
)

func Start(i models.Info, pk []byte) {
	info = i
	privateKey = pk
	server := coalaServer.NewServer(privateKey)

	server.GET("/info", getInfo)
	server.POST("/cmd", execCmd)

	log.Fatal(server.Listen("127.0.0.1:5683"))
}
