package helpers

import (
	log "device-go/shared/golog"
	"errors"
	"time"
)

func RoundRobin(cb func() error, interval time.Duration, attempts int) error {
	for {
		// if attempts == -1 then infinity
		if attempts > 0 {
			attempts--
		}
		err := cb()
		if err == nil {
			break
		}
		if err != nil {
			log.Error(err)
		}
		time.Sleep(interval)
		if attempts == 0 {
			return errors.New("round robin [max attempts]")
		}
	}
	return nil
}
