package helpers

import (
	"device-go/shared/config"
	"errors"
	log "github.com/ndmsystems/golog"
	"os"
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

func SaveContractLocal(contract []byte) error {
	f, err := os.Create(config.Get("device.contractFile"))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(contract)
	if err != nil {
		return err
	}
	return nil
}
