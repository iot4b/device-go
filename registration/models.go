package registration

// метод /register принимает на вход info текущего устройства для регистрации в блокчейне
type register struct {
	Address string `json:"a,omitempty"` // contract address if already deployed
	Version string `json:"ver"`         // версия прошивки
	Type    string `json:"t"`           // название модели устройства
	Vendor  string `json:"v"`           // вендор
}

// getEndpoints возвращает список активных нод в таком формате
type node struct {
	IpPort  string `json:"ipPort"`  // ip:port ноды
	Account string `json:"account"` // адрес смарт-контракта
}
