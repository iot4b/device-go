package dsm

// DeviceContract - контракт Vendor'a имеет возможность делать подпись транзакций для девайса
type DeviceContract struct {
	Address EverAddress `json:"address,omitempty"` //ever SC address текущего Device
	Node    EverAddress `json:"node,omitempty"`    //ever SC address Node, с которой девайс создал последнее соединение
	Elector EverAddress `json:"elector"`           //ever SC адрес Elector'a, который обслуживает сеть нод для текущего девайса
	Vendor  EverAddress `json:"vendor"`            //ever SC address производителя текущего девайса. по-умолчанию из конфигов берем

	Owners map[string]any `json:"owners"` // owners data: public_key => contract_address

	Active     bool   `json:"active"`               // if device is active
	Lock       bool   `json:"lock"`                 // if device is locked
	Stat       bool   `json:"stat"`                 // нужно ли девайсу слать статистику
	Events     bool   `json:"events"`               // sending events
	Type       string `json:"dtype,omitempty"`      // модель/тип девайса
	Version    string `json:"version,omitempty"`    // версия текущей прошивки на девайсе
	VendorName string `json:"vendorName,omitempty"` // название производителя
	VendorData string `json:"vendorData,omitempty"` // данные, которые идут от производителя девайса

	LastRegisterTime string `json:"lastRegisterTime,omitempty"` // last registration timestamp

	Hash string `json:"hash"` // hash of current contract code (contract version identifier)
}
