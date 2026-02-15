package security

import (
	"crypto/rand"     // Secure random number generator for salt
	"crypto/subtle"   // Provides constant-time compare to prevent timing attacks
	"encoding/base64" // Encoding for storing hash and salt as strings
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2" // Argon2id key derivation function
)

// ArgonParams defines parameters for Argon2id hashing
type ArgonParams struct {
	Memory      uint32 // Memory usage in KiB
	Iterations  uint32 // Number of iterations
	Parallelism uint8  // Number of threads
	SaltLength  uint32 // Length of random salt
	KeyLength   uint32 // Length of derived key (hash)
}

// DefaultParams provides reasonable defaults for password hashing
var DefaultParams = &ArgonParams{
	Memory:      64 * 1024, // 64 MB
	Iterations:  3,         // 3 iterations
	Parallelism: 2,         // Use 2 threads
	SaltLength:  16,        // 16 bytes salt
	KeyLength:   32,        // 32 bytes output hash
}

// HashPassword generates an Argon2id hash for a given password
// Returns a string in the format:
// $argon2id$v=19$m=65536,t=3,p=2$<base64_salt>$<base64_hash>
func HashPassword(password string, p *ArgonParams) (string, error) {
	// Generate random salt
	salt := make([]byte, p.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Compute Argon2id hash
	hash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	// Encode salt and hash in base64 (URL-safe, no padding)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return formatted string with all parameters, salt, and hash
	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		p.Memory,
		p.Iterations,
		p.Parallelism,
		b64Salt,
		b64Hash,
	), nil
}

// VerifyPassword checks if a plaintext password matches the stored Argon2id hash
func VerifyPassword(password, encoded string) (bool, error) {
	// Split stored hash into components
	parts := strings.SplitN(encoded, "$", 6)
	if len(parts) != 6 {
		return false, fmt.Errorf("invalid hash format")
	}

	// Validate algorithm
	if parts[1] != "argon2id" {
		return false, fmt.Errorf("unsupported algorithm: %s", parts[1])
	}

	// Validate Argon2 version
	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil || version != argon2.Version {
		return false, fmt.Errorf("incompatible argon2 version")
	}

	// Extract memory, iterations, parallelism
	var memory uint32
	var iterations uint32
	var parallelism uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism); err != nil {
		return false, err
	}

	// Decode salt and hash from base64
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	// Check minimum lengths to avoid invalid hashes
	if len(salt) < 8 || len(hash) < 16 {
		return false, fmt.Errorf("invalid salt or hash length")
	}

	// Recompute hash using the extracted parameters
	computed := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(hash)))

	// Constant-time comparison to prevent timing attacks
	if subtle.ConstantTimeCompare(hash, computed) == 1 {
		return true, nil
	}

	return false, nil
}
