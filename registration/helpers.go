package registration

import (
	"github.com/coalalib/coalago"
	log "github.com/ndmsystems/golog"
	"time"
)

func ping(nodeHost string) (duration time.Duration, err error) {
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
