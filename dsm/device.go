package dsm

// DeviceContract - контракт Vendor'a имеет возможность делать подпись транзакций для девайса
type DeviceContract struct {
	Address EverAddress `json:"address,omitempty"` //ever SC address текущего Device
	Node    EverAddress `json:"node,omitempty"`    //ever SC address Node, с которой девайс создал последнее соединение
	Elector EverAddress `json:"elector,omitempty"` //ever SC адрес Elector'a, который обслуживает сеть нод для текущего девайса
	Vendor  EverAddress `json:"vendor,omitempty"`  //ever SC address производителя текущего девайса. по-умолчанию из конфигов берем

	Owners []string `json:"owners,omitempty"` // owners public keys list

	Active     bool   `json:"active,omitempty"`     // if device is active
	Lock       bool   `json:"lock,omitempty"`       // if device is locked
	Stat       bool   `json:"stat,omitempty"`       // нужно ли девайсу слать статистику
	Type       string `json:"dtype,omitempty"`      // модель/тип девайса
	Version    string `json:"version,omitempty"`    // версия текущей прошивки на девайсе
	VendorName string `json:"vendorName,omitempty"` // название производителя
	VendorData string `json:"vendorData,omitempty"` // данные, которые идут от производителя девайса

	LastRegisterTime string `json:"lastRegisterTime,omitempty"` // last registration timestamp
}
