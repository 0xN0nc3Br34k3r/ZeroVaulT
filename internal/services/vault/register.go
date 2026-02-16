package vault

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"time"

	"github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/models"
	"github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/repository"
	vaultrepo "github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/repository/vault"

	"github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/security"
)

// Register initializes a new ZeroVaulT instance by creating the master
// authentication record inside the database. This operation is allowed
// only once; if authentication data already exists, registration is
// rejected to prevent accidental overwrites.
//
// Steps performed:
//  1. Check whether the vault is already initialized (auth bucket populated)
//  2. Hash the provided master password using Argon2id
//  3. Generate a fresh KDF salt for future key derivation
//  4. Construct an Auth model containing password hash + metadata
//  5. Serialize the Auth model to JSON
//  6. Store the record inside the "auth" bucket under key "auth_data"
//
// Returns an error if:
//   - the vault is already initialized
//   - password hashing fails
//   - random salt generation fails
//   - JSON serialization fails
//   - database insertion fails
func Register(db *repository.Database, masterPassword string) error {
	repo := vaultrepo.NewVaultRepository(db)

	// Check if auth data already exists — registration is one‑time only.
	_, err := repo.GetData("auth", "auth_data")
	if err == nil {
		return fmt.Errorf("vault already initialized; registration is not allowed")
	}

	// Derive a secure password hash using Argon2id.
	hash, err := security.HashPassword(masterPassword, security.DefaultParams)
	if err != nil {
		return err
	}

	// Generate a random salt for future key derivation.
	kdfSalt := make([]byte, 16)
	rand.Read(kdfSalt)

	auth := models.Auth{
		PasswordHash: hash,
		KDFSalt:      kdfSalt,
		KeyVersion:   1,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	authKey := "auth_data"

	// Serialize the authentication record.
	authData, err := json.Marshal(auth)
	if err != nil {
		return err
	}

	// Persist the auth record into the "auth" bucket.
	if err := repo.InsertData("auth", authKey, string(authData)); err != nil {
		return err
	}

	return nil
}
