package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	log "github.com/ndmsystems/golog"
	"strings"
	"testing"
)

func init() {
	log.Init(false)
}

func TestEncryptDecryptChaCha20Poly1305(t *testing.T) {
	// Generate new ed25519 key pair for both sides
	pub1, priv1, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}
	pub2, priv2, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Initialize keys instances
	k1 := &keys{
		Secret: hex.EncodeToString(priv1.Seed()),
	}
	k2 := &keys{
		Secret: hex.EncodeToString(priv2.Seed()),
	}

	message := "Hello, 1!"

	// Encrypt
	encrypted, err := k1.EncryptChaCha20Poly1305([]byte(message), hex.EncodeToString(pub2))
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Decrypt
	decrypted, err := k2.DecryptChaCha20Poly1305(encrypted, hex.EncodeToString(pub1))
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if decrypted != message {
		t.Errorf("Decrypted message mismatch: got %q, want %q", decrypted, message)
	}

	message = "Hello, 2!"

	// Encrypt
	encrypted, err = k2.EncryptChaCha20Poly1305([]byte(message), hex.EncodeToString(pub1))
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Decrypt
	decrypted, err = k1.DecryptChaCha20Poly1305(encrypted, hex.EncodeToString(pub2))
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if decrypted != message {
		t.Errorf("Decrypted message mismatch: got %q, want %q", decrypted, message)
	}
}

func TestDecryptWithWrongKeyFails(t *testing.T) {
	// Generate correct recipient key pair
	pubCorrect, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate correct key pair: %v", err)
	}

	// Generate wrong recipient key pair
	_, privWrong, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate wrong key pair: %v", err)
	}

	// Sender key
	senderPub, senderPriv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate sender key pair: %v", err)
	}

	// Set up keys struct for sender and wrong recipient
	sender := &keys{Secret: hex.EncodeToString(senderPriv.Seed())}
	wrongRecipient := &keys{Secret: hex.EncodeToString(privWrong.Seed())}

	message := "Top Secret"

	// Encrypt using correct recipient
	encrypted, err := sender.EncryptChaCha20Poly1305([]byte(message), hex.EncodeToString(pubCorrect))
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Try to decrypt using WRONG recipient (should fail)
	_, err = wrongRecipient.DecryptChaCha20Poly1305(encrypted, hex.EncodeToString(senderPub))
	if err == nil {
		t.Errorf("Expected decryption to fail with wrong key, but it succeeded")
	} else {
		t.Logf("Decryption correctly failed: %v", err)
	}
}

func TestEncryptWithInvalidHexKeyFails(t *testing.T) {
	k := &keys{
		Secret: strings.Repeat("a", 64), // valid dummy hex private key
	}

	_, err := k.EncryptChaCha20Poly1305([]byte("test"), "zzzzzzzz") // Invalid hex
	if err == nil {
		t.Error("Expected encryption to fail due to invalid recipient public key, but it succeeded")
	}
}

func TestDecryptWithInvalidBase64Fails(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(nil)
	k := &keys{
		Secret: hex.EncodeToString(priv.Seed()),
	}

	_, err := k.DecryptChaCha20Poly1305("!!!invalid-base64!!!", hex.EncodeToString(pub))
	if err == nil {
		t.Error("Expected decryption to fail due to invalid base64 input, but it succeeded")
	}
}
