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

// ReadSecretItem retrieves and decrypts a stored vault item from ZeroVaulT.
// This operation verifies the master password, loads the encrypted record,
// decrypts its contents using the derived vault key, and prints the resulting
// SecretItem in human‑readable JSON format.
//
// Steps performed:
//  1. Load and decode the stored authentication record
//  2. Verify the provided master password against the stored Argon2id hash
//  3. Fetch the encrypted VaultItem from the "vault_items" bucket
//  4. Deserialize the VaultItem metadata (ID, timestamps, ciphertext)
//  5. Derive the vault encryption key using Argon2id + stored KDF salt
//  6. Decrypt the VaultItem's AES‑GCM ciphertext using static AAD ("vault‑v1")
//  7. Pretty‑print the decrypted SecretItem JSON to stdout
//
// Returns an error if:
//   - authentication data cannot be loaded
//   - the master password is invalid
//   - the vault item cannot be found
//   - JSON decoding fails
//   - decryption fails
//   - pretty‑printing fails
//
// This function does not modify stored data; it only reads and decrypts.

func ReadSecretItem(db *repository.Database, masterPassword, itemKey string) error {
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

	itemRepo := vaultitem.NewVaultItemRepository(db)

	vaultData, err := itemRepo.GetData("vault_items", itemKey)
	if err != nil {
		return err
	}

	var vaultItem models.VaultItem
	if err := json.Unmarshal([]byte(vaultData), &vaultItem); err != nil {
		return err
	}

	vaultKey := security.DeriveVaultKey(masterPassword, auth.KDFSalt)

	plainText, err := security.DecryptAESGCM(vaultKey, vaultItem.Data, []byte("vault-v1"))
	if err != nil {
		return err
	}

	var pretty bytes.Buffer
	if err := json.Indent(&pretty, plainText, "", "  "); err != nil {
		panic(err)
	}

	fmt.Println(pretty.String())

	return nil
}
