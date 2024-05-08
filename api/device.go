package api

import (
	"device-go/storage"
	log "github.com/ndmsystems/golog"
	"strconv"
)

func GetBalance() float64 {
	input := map[string]string{"address": string(storage.Get().Address)}
	res, err := GET("balance", input)
	if err != nil {
		log.Error("GET(balance):", err)
		return 0
	}
	b, err := strconv.ParseFloat(string(res), 64)
	if err != nil {
		log.Error("strconv.ParseFloat:", err)
		return 0
	}
	return b
}
