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
//  key.Sign(msg) (sign string)
//  key.Validate(message, sign, pub_key) bool
//  Load(file) -> KeyPair
//  key.Public() -> string

var KeyPair keyPair

type keyPair struct {
	Public string
	Secret string
}

// Sign unsigned message using sign key pair, returns signed message
func (k *keyPair) Sign(unsigned string) string {
	res, err := everscale.Sign(unsigned, k.Public, k.Secret)
	if err != nil {
		return ""
	}
	return sign.Signature
}

func Init() {
	file, err := os.Open(config.Get("localFiles.keys"))
	defer file.Close()

	if err == nil { // get data from existing keys file
		data, err := io.ReadAll(file)
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
