package registration

import "device-go/dsm"

// метод /register принимает на вход info текущего устройства для регистрации в блокчейне
type registerRequest struct {
	Vendor     dsm.EverAddress `json:"v"`            // адрес вендора
	Key        string          `json:"k"`            // уникальный public key, который передаем для создания контракта
	Version    string          `json:"ver"`          // версия прошивки
	Type       string          `json:"t,omitempty"`  // название модели устройства
	VendorData string          `json:"vd,omitempty"` //происзолный блок данных в любом формате
}

// getEndpoints возвращает список активных нод в таком формате
type node struct {
	IpPort  string `json:"ipPort"`  // ip:port ноды
	Account string `json:"account"` // адрес смарт-контракта
}
