package cryptoKeys

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"device-go/everscale"
	"device-go/shared/config"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	log "github.com/ndmsystems/golog"
)

//TODO на замену это все
//  key.Sign(msg) (sign string)
//  key.Validate(message, sign, pub_key) bool
//  Load(file) -> KeyPair
//  key.Public() -> string

var KeyPair keyPair

type keyPair struct {
	pub  []byte
	priv ed25519.PrivateKey
}

func (k *keyPair) PublicStr() string {
	return fmt.Sprintf("%x", k.pub)
}

func (k *keyPair) SecretStr() string {
	return fmt.Sprintf("%x", k.priv)
}

func (k *keyPair) Public() []byte {
	return k.pub
}

func (k *keyPair) Secret() []byte {
	return k.priv
}

func (k *keyPair) Seed() []byte {
	return k.priv.Seed()
}

func (k *keyPair) Verify(msg []byte) bool {
	digest := sha256.Sum256(msg)
	sig := ed25519.Sign(k.priv, digest[:])
	return ed25519.Verify(k.Public(), msg, sig)
}

func (k *keyPair) setPublic(key string) {
	data, err := hex.DecodeString(key)
	if err != nil {
		log.Error(err)
		return
	}
	k.pub = data
}

func (k *keyPair) setSecret(key string) {
	data, err := hex.DecodeString(key)
	if err != nil {
		log.Error(err)
		return
	}
	k.priv = ed25519.PrivateKey{}
	k.priv = bytes.Clone(data)
}

func Init() {
	// читакм файл. если нет его, то генерим новый
	var data []byte
	keysFile, err := os.Open(config.Get("localFiles.keys"))
	if err != nil && errors.Is(err, os.ErrNotExist) {
		keys, err := everscale.GenerateKeyPair()
		if err != nil {
			log.Fatal(err)
		}
		k := keyPair{}
		k.setPublic(keys.Public)
		k.setSecret(keys.Secret)

		// пишем в файл
		log.Debugf("everscale generated keys: %+v", keys)
		log.Debugf("converted keys to ed25519: public: %s secret: %s", k.PublicStr(), k.SecretStr())
		data, err = json.Marshal(keys)
		if err != nil {
			log.Fatal(err)
		}
		err = os.WriteFile(config.Get("localFiles.keys"), data, 0644)
		if err != nil {
			log.Fatal(err)
		}
		log.Debug(string(data))
	}
	defer keysFile.Close()

	if len(data) == 0 {
		data, err = io.ReadAll(keysFile)
		if err != nil {
			log.Fatal(err)
		}
	}
	//log.Debug(data, string(data))
	KeyPair = keyPair{}
	keys := map[string]string{}
	err = json.Unmarshal(data, &keys)
	if err != nil {
		log.Fatal(err)
	}
	KeyPair.setPublic(keys["public"])
	KeyPair.setSecret(keys["secret"])
	//log.Debugf("%+v", KeyPair)
}
