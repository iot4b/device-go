package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"

	"filippo.io/edwards25519"
	log "github.com/ndmsystems/golog"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/curve25519"
)

// EncryptChaCha20Poly1305 encrypts plaintext using the recipient's Ed25519 public key.
func (k *keys) EncryptChaCha20Poly1305(
	plaintext []byte,
	recipientHexPubKey string,
) (string, error) {
	log.Debugf("EncryptChaCha20Poly1305: %s %s", plaintext, recipientHexPubKey)

	recipientPubKey, err := hex.DecodeString(recipientHexPubKey)
	if err != nil {
		return "", fmt.Errorf("invalid recipient public key: %w", err)
	}

	sharedSecret, err := k.computeSharedSecret(recipientPubKey)
	if err != nil {
		return "", fmt.Errorf("failed to compute shared secret: %w", err)
	}

	aead, err := chacha20poly1305.New(sharedSecret)
	if err != nil {
		return "", fmt.Errorf("failed to create AEAD: %w", err)
	}

	nonce := make([]byte, chacha20poly1305.NonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := aead.Seal(nil, nonce, plaintext, nil)

	combined := append(nonce, ciphertext...)
	return base64.StdEncoding.EncodeToString(combined), nil
}

// DecryptChaCha20Poly1305 decrypts base64-encoded ciphertext using the sender's Ed25519 public key.
func (k *keys) DecryptChaCha20Poly1305(
	encodedCipher string,
	senderHexPubKey string,
) (string, error) {
	log.Debug("DecryptChaCha20Poly1305:", encodedCipher, senderHexPubKey)

	senderPubKey, err := hex.DecodeString(senderHexPubKey)
	if err != nil {
		return "", fmt.Errorf("invalid sender public key: %w", err)
	}

	sharedSecret, err := k.computeSharedSecret(senderPubKey)
	if err != nil {
		return "", fmt.Errorf("failed to compute shared secret: %w", err)
	}

	encrypted, err := base64.StdEncoding.DecodeString(encodedCipher)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}
	if len(encrypted) < chacha20poly1305.NonceSize {
		return "", errors.New("ciphertext too short")
	}

	// first 12 chars of data is a nonce
	nonce := encrypted[:chacha20poly1305.NonceSize]
	ciphertext := encrypted[chacha20poly1305.NonceSize:]

	aead, err := chacha20poly1305.New(sharedSecret)
	if err != nil {
		return "", fmt.Errorf("failed to create AEAD: %w", err)
	}

	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %w", err)
	}

	return string(plaintext), nil
}

// computeSharedSecret derives a shared secret using the local private key and a remote Ed25519 public key.
func (k *keys) computeSharedSecret(remoteEdPubKey ed25519.PublicKey) ([]byte, error) {
	log.Debug("computeSharedSecret with", hex.EncodeToString(remoteEdPubKey))

	privKeyBytes, err := hex.DecodeString(k.Secret)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %w", err)
	}

	// Hash and clamp the private scalar
	h := sha512.Sum512(privKeyBytes[:32])
	scalar := h[:32]
	scalar[0] &= 248
	scalar[31] &= 127
	scalar[31] |= 64

	var compressed [32]byte
	copy(compressed[:], remoteEdPubKey[:])
	edPoint, err := new(edwards25519.Point).SetBytes(compressed[:])
	if err != nil {
		return nil, fmt.Errorf("invalid Ed25519 point: %w", err)
	}

	montPoint := edPoint.BytesMontgomery()
	sharedSecret, err := curve25519.X25519(scalar, montPoint)
	if err != nil {
		return nil, fmt.Errorf("x25519 multiplication failed: %w", err)
	}

	return sharedSecret, nil
}
