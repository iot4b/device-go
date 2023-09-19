package dsm

import (
	"crypto/ed25519"
	"time"
)

// DeviceContract - контракт Vendor'a имеет возможность делать подпись транзакций для девайса
type DeviceContract struct {
	Address       EverAddress `json:"address,omitempty"` //ever SC address текущего Device
	Node          EverAddress `json:"node,omitempty"`    //ever SC address Node, с которой девайс создал последнее соединение
	Elector       EverAddress `json:"elector,omitempty"` //ever SC адрес Elector'a, который обслуживает сеть нод для текущего девайса
	VendorAddress EverAddress `json:"vendor,omitempty"`  //ever SC address производителя текущего девайса. по-умолчанию из конфигов берем

	PublicKey []byte             `json:"publicKey,omitempty"`
	SecretKey ed25519.PrivateKey `json:"-"`

	Stat       bool   `json:"stat,omitempty"`       // нужно ли девайсу слать статистику
	Type       string `json:"type,omitempty"`       // модель/тип девайса
	Version    string `json:"version,omitempty"`    // версия текущей прошивки на девайсе
	VendorName string `json:"vendorName,omitempty"` // название производителя
	VendorData string `json:"vendorData,omitempty"` // данные, которые идут от производителя девайса
}

type CMD struct {
	Cmd   string `json:"cmd"`   // команда, которую необходимо выполнить
	Sight string `json:"sight"` // подпись, которую нужно
	Uid   string `json:"uid"`   // todo уникальный uid (для чего ???)
}

// Info - собирает данные по устройству с момента старта. Отдается при запросе на coap://device/info
// Инициируется в shared. Может быть прочитан из любого места в коде
type Info struct {
	Address EverAddress `json:"address"` // этот девайс
	Vendor  EverAddress `json:"vendor"`  // производитель
	Node    EverAddress `json:"node"`    // текущая нода
	Elector EverAddress `json:"elector"` // электор группирует устройства в сеть

	Key        string `json:"key"`                  // публичный ключ девайса. является уникальным идентификатором
	Version    string `json:"version"`              // todo версия прошивки (или модель самого девайса ???)
	Type       string `json:"type"`                 // тип девайса
	VendorName string `json:"vendorName"`           //происзолный блок данных в любом формате
	VendorData string `json:"vendorData,omitempty"` //происзолный блок данных в любом формате

	Uptime  string    `json:"uptime"`  // uptime девайса от runFrom. Обновляем при каждом чтении из Info
	RunFrom time.Time `json:"runFrom"` // время последнего запуска
}
