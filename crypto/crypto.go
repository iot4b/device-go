package crypto

import (
	"device-go/everscale"
	"encoding/json"
	"io"
	"os"

	"github.com/jinzhu/copier"

	log "github.com/ndmsystems/golog"
)

//TODO на замену это все
//  key := Generate() (KeyPair)
//  key := Load(file) (KeyPair)
//  key.Sign(msg) (sign string)
//  key.Public() (string)
//
//  Save(path, key) err
//  Validate(msg.body, message.sign, msg.sender) (bool)

var KeyPair keyPair

type keyPair struct {
	Public string `json:"public"`
	Secret string `json:"secret"`
}

// Sign unsigned message using sign key pair, returns signed message
func (k *keyPair) Sign(unsigned string) string {
	res, err := everscale.Sign(unsigned, k.Public, k.Secret)
	if err != nil {
		return ""
	}
	return res.Signed
}

// Verify signed message using public key, returns unsigned message and a flag
func (k *keyPair) Verify(signed string) (string, bool) {
	res, err := everscale.VerifySignature(signed, k.Public)
	if err != nil {
		return "", false
	}
	return res.Unsigned, true
}

// Init key storage, load from existing file or generate a new one
func Init(path string) {
	var err error
	KeyPair, err = load(path)
	if err != nil {
		KeyPair, err = generate()
		if err != nil {
			log.Fatal(err)
		}
		err = save(path, KeyPair)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// load key pair from json file
func load(path string) (k keyPair, err error) {
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

// generate everscale key pair
func generate() (k keyPair, err error) {
	keys, err := everscale.GenerateKeyPair()
	if err != nil {
		return
	}
	copier.Copy(&k, keys)

	return
}

// save key pair to json file
func save(path string, key keyPair) error {
	data, err := json.Marshal(key)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
