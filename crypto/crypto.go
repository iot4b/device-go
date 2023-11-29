package crypto

import (
	"device-go/everscale"
	"encoding/base64"
	"encoding/json"
	"github.com/markgenuine/ever-client-go/domain"
	"io"
	"os"

	log "github.com/ndmsystems/golog"
)

var Keys keys

type keys struct {
	PublicSign string `json:"public_sign"` // public key for signing
	PublicNacl string `json:"public_nacl"` // public key for nacl box encryption
	Secret     string `json:"secret"`      // shared secret key for both public keys
}

// Sign unsigned message using sign key pair, returns signed message
func (k *keys) Sign(unsigned string) string {
	res, err := everscale.Sign(unsigned, k.PublicSign, k.Secret)
	if err != nil {
		return ""
	}
	return res.Signed
}

// Verify signed message using public key, returns unsigned message and a flag
func (k *keys) Verify(signed string) (string, bool) {
	res, err := everscale.VerifySignature(signed, k.PublicSign)
	if err != nil {
		return "", false
	}
	return res.Unsigned, true
}

// Decrypt data with nacl box using sender public key.
// First 48 chars of data is a nonce hex string
func (k *keys) Decrypt(data, sender string) (string, error) {
	encrypted := data[48:]
	nonce := data[:48]

	res, err := everscale.Ever.Crypto.NaclBoxOpen(&domain.ParamsOfNaclBoxOpen{
		Encrypted:   encrypted,
		Nonce:       nonce,
		TheirPublic: sender,
		Secret:      Keys.Secret,
	})
	if err != nil {
		return "", err
	}

	decrypted, err := base64.StdEncoding.DecodeString(res.Decrypted)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

// Init key storage, load from existing file or generate a new one
func Init(path string) {
	var err error
	Keys, err = load(path)
	if err != nil {
		Keys, err = generate()
		if err != nil {
			log.Fatal(err)
		}
		err = save(path, Keys)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// load key pair from json file
func load(path string) (k keys, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &k)

	return
}

// generate everscale crypto keys
func generate() (k keys, err error) {
	sign, err := everscale.Ever.Crypto.GenerateRandomSignKeys()
	if err != nil {
		return
	}
	nacl, err := everscale.Ever.Crypto.NaclBoxKeypairFromSecretKey(
		&domain.ParamsOfNaclBoxKeyPairFromSecret{Secret: sign.Secret},
	)
	if err != nil {
		return
	}

	k.Secret = sign.Secret
	k.PublicSign = sign.Public
	k.PublicNacl = nacl.Public

	return
}

// save key pair to json file
func save(path string, key keys) error {
	data, err := json.Marshal(key)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
