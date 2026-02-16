package security

import (
	"golang.org/x/crypto/argon2"
)

// DeriveVaultKey derives a symmetric encryption key from a user-provided
// password and a cryptographically secure salt using Argon2id.
//
// Argon2id is chosen because it provides strong resistance against both
// GPU‑based brute‑force attacks (memory‑hard) and side‑channel attacks.
//
// Parameters (Iterations, Memory, Parallelism, KeyLength) are taken from
// DefaultParams, allowing centralized tuning of security settings.
//
// The returned key is deterministic for the same (password, salt) pair,
// and should be used as the master key for encrypting the vault.
func DeriveVaultKey(password string, salt []byte) []byte {
	return argon2.IDKey(
		[]byte(password),
		salt,
		DefaultParams.Iterations,
		DefaultParams.Memory,
		DefaultParams.Parallelism,
		DefaultParams.KeyLength,
	)
}
