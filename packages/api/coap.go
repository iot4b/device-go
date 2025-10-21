package api

import (
	"device-go/packages/config"
	"device-go/packages/storage"
	"encoding/json"

	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
	"github.com/pkg/errors"
)

var client = coalago.NewClient()

func GET(path string, data any) ([]byte, error) {
	msg, err := build(coalago.GET, path, data)
	if err != nil {
		return nil, err
	}
	return send(msg, path, "GET")
}

func POST(path string, data any) ([]byte, error) {
	msg, err := build(coalago.POST, path, data)
	if err != nil {
		return nil, err
	}
	return send(msg, path, "POST")
}

func build(method coalago.CoapCode, path string, data any) (*coalago.CoAPMessage, error) {
	msg := coalago.NewCoAPMessage(coalago.CON, method)
	msg.SetURIPath(path)
	msg.Timeout = config.Time("timeout.coala")

	if data != nil {
		if method == coalago.GET {
			for k, v := range data.(map[string]string) {
				msg.SetURIQuery(k, v)
			}
		}
		if method == coalago.POST {
			var b []byte
			var ok bool
			var err error
			if b, ok = data.([]byte); !ok {
				b, err = json.Marshal(data)
				if err != nil {
					return nil, err
				}
			}
			msg.SetStringPayload(string(b))
		}
	}

	return msg, nil
}

func send(msg *coalago.CoAPMessage, path, method string) ([]byte, error) {
	log.Debug("coap:", method, storage.Device.NodeIpPort, path, msg.Payload.String())
	resp, err := client.Send(msg, storage.Device.NodeIpPort)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("nil response")
	}
	if resp.Code != coalago.CoapCodeContent {
		return nil, errors.New(string(resp.Body))
	}
	return resp.Body, nil
}
