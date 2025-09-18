package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"device-go/packages/utils"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/nacl/sign"

	log "github.com/ndmsystems/golog"
)

var (
	Keys     keys
	filePath string
)

type keys struct {
	PublicSign string `json:"public_sign"` // public key for signing
	PublicNacl string `json:"public_nacl"` // public key for nacl box encryption
	Secret     string `json:"secret"`      // shared secret key for both public keys
}

// Sign unsigned message using sign key pair, returns signature
func (k *keys) Sign(unsigned []byte) []byte {
	private, err := hex.DecodeString(k.Secret)
	if err != nil {
		return nil
	}
	public, err := hex.DecodeString(k.PublicSign)
	if err != nil {
		return nil
	}
	signature := sign.Sign(nil, unsigned, (*[64]byte)(append(private, public...)))

	return signature[:64]
}

// VerifySignature reports whether sig is a valid signature of message by public key
func VerifySignature(pubKey string, message []byte, sig string) bool {
	public, err := hex.DecodeString(pubKey)
	if err != nil {
		return false
	}
	sigDecoded, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return false
	}

	return ed25519.Verify(public, message, sigDecoded)
}

// Decrypt data with nacl box using sender public key.
func (k *keys) Decrypt(data, sender string) (string, error) {
	encrypted, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	// first 24 chars of data is a nonce
	var nonce [24]byte
	copy(nonce[:], encrypted[:24])

	decrypted, ok := box.Open(nil, encrypted[24:], &nonce, hexTo32b(sender), hexTo32b(Keys.Secret))
	if !ok {
		return "", errors.New("decryption error")
	}

	return string(decrypted), nil
}

// Init key storage, load from existing file or generate a new one
func Init(fileName string) {
	filePath = filepath.Join(utils.GetFilesDir(), fileName)

	var err error
	Keys, err = load()
	if err != nil {
		Keys, err = generate()
		if err != nil {
			log.Fatal(err)
		}
		err = save()
		if err != nil {
			log.Fatal(err)
		}
	}
}

// load key pair from json file
func load() (k keys, err error) {
	file, err := os.Open(filePath)
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

// generate crypto keys for signing and nacl box encryption
func generate() (k keys, err error) {
	// generate key pair for signing
	publicSign, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return
	}

	// use private key to calculate public key for nacl box encryption
	publicNacl := new([32]byte)
	private32b := new([32]byte)
	copy(private32b[:], private[:32])
	curve25519.ScalarBaseMult(publicNacl, private32b)

	k.Secret = hex.EncodeToString(private32b[:])
	k.PublicSign = hex.EncodeToString(publicSign)
	k.PublicNacl = hex.EncodeToString(publicNacl[:])

	return
}

// save key pair to json file
func save() error {
	data, err := json.Marshal(Keys)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

// hexTo32b converts hex string to 32-byte array
// string length must be 64 chars (32 hexadecimal pairs)
func hexTo32b(h string) *[32]byte {
	b, err := hex.DecodeString(h)
	if err != nil {
		return nil
	}
	var res [32]byte
	copy(res[:], b)
	return &res
}
