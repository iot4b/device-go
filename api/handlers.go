package api

import (
	"github.com/coalalib/coalago/message"
	"github.com/coalalib/coalago/resource"
)

type info struct {
	Key     string `json:"key"`
	Version string `json:"version"`
	Type    string `json:"type"`
	Vendor  string `json:"vendor"`
}

type cmd struct {
	Cmd   string `json:"cmd"`
	Sight string `json:"sight"`
	Uid   string `json:"uid"`
}

func getInfo(message *coalaMsg.CoAPMessage) *resource.CoAPResourceHandlerResult {
	return nil
}

func execCmd(message *coalaMsg.CoAPMessage) *resource.CoAPResourceHandlerResult {
	return nil
}
