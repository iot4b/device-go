package registration

type register struct {
	Key     string `json:"key"`     // public key для подписи сообщений
	Version string `json:"version"` // версия прошивки
	Type    string `json:"type"`    // название модели устройства
	Vendor  string `json:"vendor"`  // вендор
}
