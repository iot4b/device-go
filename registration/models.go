package registration

import "device-go/dsm"

// метод /register принимает на вход info текущего устройства для регистрации в блокчейне
type registerRequest struct {
	Address    dsm.EverAddress `json:"a,omitempty"`  // device contract address if deployed
	Elector    dsm.EverAddress `json:"e"`            // elector address
	Vendor     dsm.EverAddress `json:"v"`            // адрес вендора
	Owners     []string        `json:"o"`            // owners public keys list
	PublicKey  string          `json:"k"`            // уникальный public key, который передаем для создания контракта
	Version    string          `json:"ver"`          // версия прошивки
	Type       string          `json:"t,omitempty"`  // название модели устройства
	VendorName string          `json:"vn,omitempty"` //происзолный блок данных в любом формате
	VendorData string          `json:"vd,omitempty"` //происзолный блок данных в любом формате
}

// getEndpoints возвращает список активных нод в таком формате
type node struct {
	IpPort  string `json:"ip_port"` // ip:port ноды
	Account string `json:"account"` // адрес смарт-контракта
}
