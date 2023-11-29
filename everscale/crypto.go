package everscale

import "github.com/markgenuine/ever-client-go/domain"

// Sign unsigned message with public and secret keys
func Sign(unsigned, public, secret string) (*domain.ResultOfSign, error) {
	return Ever.Crypto.Sign(&domain.ParamsOfSign{
		Unsigned: unsigned,
		Keys: &domain.KeyPair{
			Public: public,
			Secret: secret,
		},
	})
}

// VerifySignature of signed message using public key, returns unsigned message
func VerifySignature(signed, public string) (*domain.ResultOfVerifySignature, error) {
	return Ever.Crypto.VerifySignature(&domain.ParamsOfVerifySignature{
		Signed: signed,
		Public: public,
	})
}
