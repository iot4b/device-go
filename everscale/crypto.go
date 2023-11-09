package everscale

import "github.com/markgenuine/ever-client-go/domain"

// GenerateKeyPair generates random key pair for signing messages in everscale
func GenerateKeyPair() (*domain.KeyPair, error) {
	return ever.Crypto.GenerateRandomSignKeys()
}

// Sign unsigned message with public and secret keys
func Sign(unsigned, public, secret string) (*domain.ResultOfSign, error) {
	return ever.Crypto.Sign(&domain.ParamsOfSign{
		Unsigned: unsigned,
		Keys: &domain.KeyPair{
			Public: public,
			Secret: secret,
		},
	})
}

// VerifySignature of signed message using public key, returns unsigned message
func VerifySignature(signed, public string) (*domain.ResultOfVerifySignature, error) {
	return ever.Crypto.VerifySignature(&domain.ParamsOfVerifySignature{
		Signed: signed,
		Public: public,
	})
}
