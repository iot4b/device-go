package crypto

import (
	"device-go/everscale"
	"device-go/shared/config"
	"encoding/json"
	"io"
	"os"

	"github.com/jinzhu/copier"

	log "github.com/ndmsystems/golog"
)

//  key := Generate() (KeyPair)
//  key := Load(file) (KeyPair)
//  key.Save(file) (err)
//  key.Sign(msg) (sign string)
//  Validate(msg.body, message.sign, msg.sender) (bool)

var KeyPair keyPair

type keyPair struct {
	Public string
	Secret string
}

func (k *keyPair) Sign(unsigned string) string {
	sign, err := everscale.Sign(unsigned, k.Public, k.Secret)
	if err != nil {
		return ""
	}
	return sign.Signature
}

func Init(path string) {
	file, err := os.Open(path)
	defer file.Close()

	if err == nil { // get data from existing keys file
		data, err := io.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(data, &KeyPair)
		if err != nil {
			log.Fatal(err)
		}
	} else { // generate new key pair and save to file
		keys, err := everscale.GenerateKeyPair()
		if err != nil {
			log.Fatal(err)
		}

		log.Debugf("everscale generated keys: %+v", keys)
		data, err := json.Marshal(keys)
		if err != nil {
			log.Fatal(err)
		}
		err = os.WriteFile(config.Get("localFiles.keys"), data, 0644)
		if err != nil {
			log.Fatal(err)
		}

		copier.Copy(&KeyPair, keys)
	}
}
