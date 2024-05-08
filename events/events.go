package events

import (
	"device-go/api"
	"device-go/dsm"
	"device-go/storage"
	log "github.com/ndmsystems/golog"
	"reflect"
)

type event interface {
	process()
}

type payload struct {
	Address dsm.EverAddress `json:"address"`
	Name    string          `json:"name"`
	Data    any             `json:"data"`
}

func Send(e event) {
	e.process()

	_, err := api.POST("device/event", payload{
		Address: storage.Get().Address,
		Name:    reflect.TypeOf(e).Elem().Name(),
		Data:    e,
	})
	if err != nil {
		log.Errorf("coap.POST: %v", err)
	}
}

type Register struct {
	Name    string  `json:"name"`
	Msg     string  `json:"msg"`
	Balance float64 `json:"balance"`
}

func (event *Register) process() {
	event.Name = "registration"
	event.Msg = "Successful"
	event.Balance = api.GetBalance()
}
