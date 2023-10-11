package everscale

import (
	log "github.com/ndmsystems/golog"

	"github.com/markgenuine/ever-client-go"
)

var ever *goever.Ever

func Init(endpoints []string) {
	var err error
	ever, err = goever.NewEver("", endpoints, "")
	if err != nil {
		log.Fatal(err)
	}
}
