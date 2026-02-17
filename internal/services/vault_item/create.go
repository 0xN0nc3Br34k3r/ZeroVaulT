package vaultitem

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/models"
	"github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/repository"
	"github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/repository/vault"
	vaultitem "github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/repository/vault_item"
	"github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/security"
)

// Create adds a new encrypted vault item to ZeroVaulT. This operation
// verifies the master password, encrypts the secret fields, and stores
// the resulting VaultItem record inside the "vault_items" bucket.
//
// Steps performed:
//  1. Load and decode the stored authentication record
//  2. Verify the provided master password against the stored Argon2id hash
//  3. Validate required fields (e.g., itemName)
//  4. Parse and normalize tags (comma‑separated string → []string)
//  5. Build a SecretItem containing all sensitive fields
//  6. Serialize the SecretItem to JSON
//  7. Derive the vault encryption key using Argon2id + stored KDF salt
//  8. Encrypt the SecretItem using AES‑GCM with a static AAD ("vault‑v1")
//  9. Wrap the ciphertext in a VaultItem model with metadata
//  10. Serialize and store the VaultItem in the "vault_items" bucket
//
// Returns an error if:
//   - authentication data cannot be loaded
//   - the master password is invalid
//   - required fields are missing (e.g., itemName)
//   - JSON serialization fails
//   - encryption fails
//   - database insertion fails
func CreateSecretItem(
	db *repository.Database,
	masterPassword, itemName, serviceName, website,
	username, password, email, notes, tags, totpSecret string,
) error {

	// Use correct repo for auth
	authRepo := vault.NewVaultRepository(db)

	data, err := authRepo.GetData("auth", "auth_data")
	if err != nil {
		return err
	}

	var auth models.Auth
	if err := json.Unmarshal([]byte(data), &auth); err != nil {
		return err
	}

	// Verify password
	ok, err := security.VerifyPassword(masterPassword, auth.PasswordHash)
	if err != nil || !ok {
		return fmt.Errorf("invalid password")
	}

	if itemName == "" {
		return fmt.Errorf("itemName cannot be empty")
	}

	// Parse tags
	rawTags := strings.Split(tags, ",")
	tagList := make([]string, 0, len(rawTags))

	for _, t := range rawTags {
		trimmed := strings.TrimSpace(t)
		if trimmed != "" {
			tagList = append(tagList, trimmed)
		}
	}

	now := time.Now()

	secretItem := &models.SecretItem{
		ServiceName: serviceName,
		Website:     website,
		Username:    username,
		Password:    password,
		Email:       email,
		Notes:       notes,
		Tags:        tagList,
		TOTPSecret:  totpSecret,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	jsonSecretItem, err := json.Marshal(secretItem)
	if err != nil {
		return err
	}

	vaultKey := security.DeriveVaultKey(masterPassword, auth.KDFSalt)

	cipherText, err := security.EncryptAESGCM(vaultKey, jsonSecretItem, []byte("vault-v1"))
	if err != nil {
		return err
	}

	vaultItem := &models.VaultItem{
		ID:        itemName,
		Data:      cipherText,
		CreatedAt: now,
		UpdatedAt: now,
	}

	jsonVaultItem, err := json.Marshal(vaultItem)
	if err != nil {
		return err
	}

	// Use correct repo for vault items
	itemRepo := vaultitem.NewVaultItemRepository(db)
	return itemRepo.InsertData("vault_items", itemName, string(jsonVaultItem))
}
