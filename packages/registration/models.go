package registration

import "device-go/packages/dsm"

// метод /register принимает на вход info текущего устройства для регистрации в блокчейне
type registerRequest struct {
	Name       string          `json:"nm,omitempty"` // device name
	Address    dsm.EverAddress `json:"a,omitempty"`  // device contract address if deployed
	Group      dsm.EverAddress `json:"g,omitempty"`  // device group contract address if any
	Elector    dsm.EverAddress `json:"e"`            // elector address
	Vendor     dsm.EverAddress `json:"v"`            // адрес вендора
	DeviceAPI  dsm.EverAddress `json:"api"`          // device API contract address
	Owners     map[string]any  `json:"o"`            // owners data: public_key => contract_address
	PublicSign string          `json:"k"`            // уникальный public key, который передаем для создания контракта
	PublicNacl string          `json:"n"`            // device public key for nacl box encryption
	Version    string          `json:"ver"`          // версия прошивки
	Type       string          `json:"t,omitempty"`  // название модели устройства
	VendorName string          `json:"vn,omitempty"` // vendor name
	VendorData string          `json:"vd,omitempty"` // произволный блок данных в любом формате
	Stat       bool            `json:"st"`           // storing statistics
	Events     bool            `json:"ev"`           // sending events
	Hash       string          `json:"h"`            // hash of current contract code (contract version identifier)
}

type registerResponse struct {
	Address dsm.EverAddress `json:"a,omitempty"` //ever SC address текущего Device
	Node    dsm.EverAddress `json:"n,omitempty"` //ever SC address Node, с которой девайс создал последнее соединение
	Elector dsm.EverAddress `json:"e,omitempty"` //ever SC адрес Elector'a, который обслуживает сеть нод для текущего девайса
	Vendor  dsm.EverAddress `json:"v,omitempty"` //ever SC address производителя текущего девайса

	Stat   bool `json:"st,omitempty"` // storing statistics
	Events bool `json:"ev,omitempty"` // sending events

	Hash string `json:"h,omitempty"` // actual contract code hash
}

// getEndpoints возвращает список активных нод в таком формате
type node struct {
	IpPort  string `json:"ip_port"` // ip:port ноды
	Account string `json:"account"` // адрес смарт-контракта
}
