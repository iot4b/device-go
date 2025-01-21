package events

import (
	"device-go/packages/api"
	"device-go/packages/dsm"
	"device-go/packages/storage"
	"reflect"

	log "github.com/ndmsystems/golog"
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
		Address: storage.Device.Address,
		Name:    reflect.TypeOf(e).Elem().Name(),
		Data:    e,
	})
	if err != nil {
		log.Errorf("coap.POST: %v", err)
	}
}

type Start struct {
	Name    string  `json:"name"`
	Msg     string  `json:"msg"`
	Balance float64 `json:"balance"`
}

func (event *Start) process() {
	event.Name = "start"
	event.Msg = string("Device started: " + storage.Device.Address)
	event.Balance = api.GetBalance()
}

type Register struct {
	Name    string  `json:"name"`
	Msg     string  `json:"msg"`
	Balance float64 `json:"balance"`
}

func (event *Register) process() {
	event.Name = "registration"
	event.Msg = string("Successful registration: " + storage.Device.Address)
	event.Balance = api.GetBalance()
}
