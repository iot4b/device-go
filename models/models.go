package models

type Info struct {
	Key     string `json:"key"`
	Version string `json:"version"`
	Type    string `json:"type"`
	Vendor  string `json:"vendor"`
}

type CMD struct {
	Cmd   string `json:"cmd"`
	Sight string `json:"sight"`
	Uid   string `json:"uid"`
}
