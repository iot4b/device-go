package models

type Info struct {
	Key     string `json:"key"`     // публичный ключ девайса. является уникальным идентификатором
	Version string `json:"version"` // todo версия прошивки (или модель самого девайса ???)
	Type    string `json:"type"`    // тип девайса
	Vendor  string `json:"vendor"`  // производитель
}

type CMD struct {
	Cmd   string `json:"cmd"`   // команда, которую необходимо выполнить
	Sight string `json:"sight"` // подпись, которую нужно
	Uid   string `json:"uid"`   // todo уникальный uid (для чего ???)
}
