package cryptoKeys

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"device-go/shared/config"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/ndmsystems/golog"

	"io"
	"os"
)

var KeyPair keyPair

type keyPair struct {
	public  []byte
	private ed25519.PrivateKey
}

func (k *keyPair) PublicStr() string {
	return fmt.Sprintf("%x", k.public)
}

func (k *keyPair) SecretStr() string {
	return fmt.Sprintf("%x", k.private)
}

func (k *keyPair) Public() []byte {
	return k.public
}

func (k *keyPair) Secret() []byte {
	return k.private
}

func (k *keyPair) Seed() []byte {
	return k.private.Seed()
}

func (k *keyPair) Verify(msg []byte) bool {
	digest := sha256.Sum256(msg)
	sig := ed25519.Sign(k.private, digest[:])
	return ed25519.Verify(k.Public(), msg, sig)
}

func (k *keyPair) setPublic(key string) {
	data, err := hex.DecodeString(key)
	if err != nil {
		log.Error(err)
		return
	}
	k.public = data
}

func (k *keyPair) setSecret(key string) {
	data, err := hex.DecodeString(key)
	if err != nil {
		log.Error(err)
		return
	}
	k.private = ed25519.PrivateKey{}
	k.private = bytes.Clone(data)
}

func Init() {
	log.Debug("init public/secret keys")
	// читакм файл. если нет его, то генерим новый
	var data []byte
	keysFile, err := os.Open(config.Get("localFiles.keys"))
	if err != nil && errors.Is(err, os.ErrNotExist) {
		public, private, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			log.Fatal(err)
		}
		k := keyPair{
			public:  public,
			private: private,
		}
		// пишем в файл
		keys := map[string]string{
			"public": fmt.Sprintf("%x", k.public),
			"secret": fmt.Sprintf("%x", k.private),
		}
		log.Debug(keys, public)
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
	log.Debug(data, string(data))
	KeyPair = keyPair{}
	keys := map[string]string{}
	err = json.Unmarshal(data, &keys)
	if err != nil {
		log.Fatal(err)
	}
	KeyPair.setPublic(keys["public"])
	KeyPair.setSecret(keys["secret"])
	log.Debugf("keypair: %+v", KeyPair)
}
