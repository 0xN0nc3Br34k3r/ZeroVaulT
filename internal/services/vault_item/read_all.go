package vaultitem

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/models"
	"github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/repository"
	"github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/repository/vault"
	vaultitem "github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/repository/vault_item"
	"github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/security"
)

// ReadAllSecretItem retrieves, decrypts, and prints every stored secret item
// inside the "vault_items" bucket. This operation verifies the master password,
// loads all encrypted VaultItem records, decrypts each one using the derived
// vault key, and prints the resulting SecretItem objects in formatted JSON.
//
// Steps performed:
//  1. Load and decode the stored authentication record
//  2. Verify the provided master password against the stored Argon2id hash
//  3. Fetch all encrypted VaultItem entries from the "vault_items" bucket
//  4. Derive the vault encryption key using Argon2id + stored KDF salt
//  5. For each item:
//     - Deserialize the VaultItem metadata
//     - Decrypt the AES‑GCM ciphertext using static AAD ("vault‑v1")
//     - Pretty‑print the decrypted SecretItem JSON
//     - Continue processing even if individual items fail
//
// This function is tolerant of partial failures: if one item cannot be decoded
// or decrypted, the error is logged and processing continues for the rest.
//
// Returns an error only if:
//   - authentication data cannot be loaded
//   - the master password is invalid
//   - fetching the list of vault items fails
//
// This function does not modify stored data; it only reads and decrypts.
func ReadAllSecretItem(db *repository.Database, masterPassword string) error {
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

	// Load all encrypted vault items
	itemRepo := vaultitem.NewVaultItemRepository(db)
	secretData, err := itemRepo.GetAllData("vault_items")
	if err != nil {
		return err
	}

	// Derive vault key for decrypting all items
	vaultKey := security.DeriveVaultKey(masterPassword, auth.KDFSalt)

	// Iterate through all stored items
	for key, raw := range secretData {

		// Decode VaultItem metadata
		var vaultItem models.VaultItem
		if err := json.Unmarshal([]byte(raw), &vaultItem); err != nil {
			fmt.Printf("FAILED to decode item %s: %v\n", key, err)
			continue
		}

		// Decrypt the stored ciphertext
		plainText, err := security.DecryptAESGCM(vaultKey, vaultItem.Data, []byte("vault-v1"))
		if err != nil {
			fmt.Printf("FAILED to decrypt key %s: %v\n", key, err)
			continue
		}

		// Pretty‑print the decrypted JSON
		var pretty bytes.Buffer
		if err := json.Indent(&pretty, plainText, "", "  "); err != nil {
			fmt.Printf("FAILED to pretty-print key %s: %v\n", key, err)
			continue
		}

		fmt.Printf("Key: %s\nValue:\n%s\n\n", key, pretty.String())
	}

	return nil
}
