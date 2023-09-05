package registration

// метод /register принимает на вход info текущего устройства для регистрации в блокчейне
type register struct {
	Key     string `json:"key"`     // public key для подписи сообщений
	Version string `json:"version"` // версия прошивки
	Type    string `json:"type"`    // название модели устройства
	Vendor  string `json:"vendor"`  // вендор
}

// getEndpoints возвращает список активных нод в таком формате
type node struct {
	IpPort  string `json:"ipPort"`  // ip:port ноды
	Account string `json:"account"` // адрес смарт-контракта
}
