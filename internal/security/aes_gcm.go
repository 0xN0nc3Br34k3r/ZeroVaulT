package security

import (
	"crypto/aes"    // AES block cipher
	"crypto/cipher" // Provides GCM and other cipher modes
	"crypto/rand"   // Secure random number generation
	"fmt"           // Formatting for errors
	"io"            // For io.ReadFull
)

// EncryptAESGCM encrypts plainText using AES in Galois/Counter Mode (GCM).
// key: AES key (16, 24, 32 bytes for AES-128/192/256)
// plainText: data to encrypt
// aad: additional authenticated data (optional, can be nil)
// Returns: encrypted data with nonce prepended
func EncryptAESGCM(key, plainText, aad []byte) ([]byte, error) {
	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Wrap block cipher in Galois/Counter Mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Generate a random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Encrypt plaintext and append authentication tag
	ciphertext := gcm.Seal(nil, nonce, plainText, aad)

	// Prepend nonce to ciphertext for use in decryption
	out := make([]byte, 0, len(nonce)+len(ciphertext))
	out = append(out, nonce...)
	out = append(out, ciphertext...)

	return out, nil
}

// DecryptAESGCM decrypts data encrypted with EncryptAESGCM.
// key: AES key used for encryption
// data: nonce + ciphertext returned from EncryptAESGCM
// aad: additional authenticated data (must match encryption)
// Returns: decrypted plaintext
func DecryptAESGCM(key, data, aad []byte) ([]byte, error) {
	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Wrap block cipher in Galois/Counter Mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	// Split data into nonce and ciphertext
	nonce := data[:nonceSize]
	gcmCiphertext := data[nonceSize:]

	// Decrypt and verify authentication tag
	plaintext, err := gcm.Open(nil, nonce, gcmCiphertext, aad)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
