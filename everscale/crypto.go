package everscale

import "github.com/markgenuine/ever-client-go/domain"

func GenerateKeyPair() (*domain.KeyPair, error) {
	keys, err := ever.Crypto.GenerateRandomSignKeys()
	return keys, err
}

func Sign(unsigned, public, secret string) (*domain.ResultOfSign, error) {
	keys := &domain.KeyPair{
		Public: public,
		Secret: secret,
	}
	return ever.Crypto.Sign(&domain.ParamsOfSign{
		Unsigned: unsigned,
		Keys:     keys,
	})
}
