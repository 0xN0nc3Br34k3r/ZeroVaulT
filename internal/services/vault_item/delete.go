package vaultitem

import (
	"encoding/json"
	"fmt"

	"github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/models"
	"github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/repository"
	"github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/repository/vault"
	vaultitem "github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/repository/vault_item"
	"github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/security"
)

// Delete removes a stored vault item from ZeroVaulT.
// This operation verifies the master password, loads the encrypted
// authentication record, confirms access rights, and deletes the
// specified VaultItem from the "vault_items" bucket.
//
// Steps performed:
//  1. Validate that the provided item key is not empty
//  2. Load and decode the stored authentication record
//  3. Verify the provided master password against the stored Argon2id hash
//  4. Initialize the VaultItem repository
//  5. Delete the corresponding VaultItem entry from the "vault_items" bucket
//
// Returns an error if:
//   - the item key is empty
//   - authentication data cannot be loaded
//   - the master password is invalid
//   - the vault item does not exist
//   - deletion fails
//
// This function permanently removes the stored item and cannot be undone.
func Delete(db *repository.Database, masterPassword, itemKey string) error {
	if itemKey == "" {
		return fmt.Errorf("key is not empty")
	}

	// Load authentication metadata
	authRepo := vault.NewVaultRepository(db)
	data, err := authRepo.GetData("auth", "auth_data")
	if err != nil {
		return err
	}

	// Decode auth record
	var auth models.Auth
	if err := json.Unmarshal([]byte(data), &auth); err != nil {
		return err
	}

	// Verify master password
	ok, err := security.VerifyPassword(masterPassword, auth.PasswordHash)
	if err != nil || !ok {
		return fmt.Errorf("invalid password")
	}

	itemRepo := vaultitem.NewVaultItemRepository(db)

	if err := itemRepo.Delete("vault_items", itemKey); err != nil {
		return err
	}

	return nil
}
