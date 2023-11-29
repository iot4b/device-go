package everscale

import (
	log "github.com/ndmsystems/golog"

	"github.com/markgenuine/ever-client-go"
)

var Ever *goever.Ever

func Init(endpoints []string) {
	var err error
	Ever, err = goever.NewEver("", endpoints, "")
	if err != nil {
		log.Fatal(err)
	}
}

// Destroy client when finished
func Destroy() {
	Ever.Client.Destroy()
}
