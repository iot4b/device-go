package registration

import (
	"device-go/helpers"
	"encoding/json"
	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
	"time"
)

type register struct {
	Key     string `json:"key"`
	Version string `json:"version"`
	Type    string `json:"type"`
	Vendor  string `json:"vendor"`
}

func Register(nodeHost, public, version, Type, vendor string) error {
	client := coalago.NewClient()

	payload := register{
		Key:     public,
		Version: version,
		Type:    Type,
		Vendor:  vendor,
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := coalago.NewCoAPMessage(coalago.CON, coalago.POST)
	msg.SetURIPath("/register")
	msg.SetStringPayload(string(bytes))
	resp, err := client.Send(msg, nodeHost)
	if err != nil {
		log.Error(err)
		return err
	}

	// сохраняем контракт девайса локально
	return helpers.SaveContractLocal(resp.Body)
}

func NodeList(masterNode string) (list []string, err error) {
	client := coalago.NewClient()

	msg := coalago.NewCoAPMessage(coalago.CON, coalago.GET)
	msg.SetURIPath("/endpoints")
	resp, err := client.Send(msg, masterNode)
	if err != nil {
		return nil, err
	}
	log.Debug(string(resp.Body))
	err = json.Unmarshal(resp.Body, &list)
	if err != nil {
		return nil, err
	}
	return
}

func Ping(nodeHost string) (duration time.Duration, err error) {
	start := time.Now()

	client := coalago.NewClient()

	msg := coalago.NewCoAPMessage(coalago.CON, coalago.GET)
	msg.SetURIPath("/info")
	resp, err := client.Send(msg, nodeHost)
	if err != nil {
		return
	}
	duration = time.Since(start)
	log.Debugf("node: %s, ping time: %d ms, %s", nodeHost, duration.Milliseconds(), string(resp.Body))
	return
}
