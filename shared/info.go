package shared

import (
	"device-go/dsm"
	"time"
)

var Info dsm.Info

func init() {
	Info = dsm.Info{RunFrom: time.Now()}
}
