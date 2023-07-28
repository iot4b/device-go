package workers

import (
	"device-go/client"
	"time"
)

func RunAlive(client *client.Client) {
	for {
		client.SendAlive()
		time.Sleep(30 * time.Second)
	}
}
