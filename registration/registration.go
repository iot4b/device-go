package registration

import (
	"encoding/json"
	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
	"time"
)

func Register(nodeHost string) error {
	client := coalago.NewClient()

	msg := coalago.NewCoAPMessage(coalago.CON, coalago.POST)
	msg.SetURIPath("/register")
	_, err := client.Send(msg, nodeHost)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
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
