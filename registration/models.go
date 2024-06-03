package registration

import "device-go/dsm"

// метод /register принимает на вход info текущего устройства для регистрации в блокчейне
type registerRequest struct {
	Address    dsm.EverAddress `json:"a,omitempty"`  // device contract address if deployed
	Group      dsm.EverAddress `json:"g,omitempty"`  // device group contract address if any
	Elector    dsm.EverAddress `json:"e"`            // elector address
	Vendor     dsm.EverAddress `json:"v"`            // адрес вендора
	Owners     map[string]any  `json:"o"`            // owners data: public_key => contract_address
	PublicSign string          `json:"k"`            // уникальный public key, который передаем для создания контракта
	PublicNacl string          `json:"n"`            // device public key for nacl box encryption
	Version    string          `json:"ver"`          // версия прошивки
	Type       string          `json:"t,omitempty"`  // название модели устройства
	VendorName string          `json:"vn,omitempty"` // vendor name
	VendorData string          `json:"vd,omitempty"` // произволный блок данных в любом формате
	Hash       string          `json:"h"`            // hash of current contract code (contract version identifier)
}

// getEndpoints возвращает список активных нод в таком формате
type node struct {
	IpPort  string `json:"ip_port"` // ip:port ноды
	Account string `json:"account"` // адрес смарт-контракта
}
