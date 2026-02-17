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

// UpdateSecretItem loads, decrypts, modifies, and re‑encrypts an existing
// secret item stored in ZeroVaulT. This operation verifies the master
// password, retrieves the encrypted VaultItem, applies only the fields
// explicitly provided by the user, and writes the updated record back
// into the "vault_items" bucket.
//
// Steps performed:
//  1. Load and decode the stored authentication record
//  2. Verify the provided master password using the stored Argon2id hash
//  3. Fetch the encrypted VaultItem from the "vault_items" bucket
//  4. Deserialize the VaultItem metadata (ID, timestamps, ciphertext)
//  5. Derive the vault encryption key using Argon2id + stored KDF salt
//  6. Decrypt the VaultItem's AES‑GCM ciphertext using static AAD ("vault‑v1")
//  7. Decode the decrypted SecretItem JSON
//  8. Apply updates only to fields explicitly provided by the caller
//  9. Update timestamps (SecretItem.UpdatedAt and VaultItem.UpdatedAt)
//
// 10. Re‑encrypt the modified SecretItem using AES‑GCM
// 11. Serialize the updated VaultItem and write it back to BoltDB
//
// Returns an error if:
//   - authentication data cannot be loaded
//   - the master password is invalid
//   - the target vault item does not exist
//   - JSON decoding or encoding fails
//   - decryption or encryption fails
//   - the updated record cannot be written to the database
//
// This function performs an in‑place update of an existing secret item.
// Fields not explicitly provided remain unchanged.
func UpdateSecretItem(
	db *repository.Database,
	masterPassword, itemName, serviceName, website,
	username, password, email, notes, tags, totpSecret string,
) error {

	authRepo := vault.NewVaultRepository(db)

	data, err := authRepo.GetData("auth", "auth_data")
	if err != nil {
		return err
	}

	var auth models.Auth
	if err := json.Unmarshal([]byte(data), &auth); err != nil {
		return err
	}

	ok, err := security.VerifyPassword(masterPassword, auth.PasswordHash)
	if err != nil || !ok {
		return fmt.Errorf("invalid password")
	}

	itemRepo := vaultitem.NewVaultItemRepository(db)

	vaultData, err := itemRepo.GetData("vault_items", itemName)
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

	var secretItems models.SecretItem
	if err := json.Unmarshal(plainText, &secretItems); err != nil {
		return err
	}

	// Apply updates only if provided
	if serviceName != "" {
		secretItems.ServiceName = serviceName
	}
	if website != "" {
		secretItems.Website = website
	}
	if username != "" {
		secretItems.Username = username
	}
	if password != "" {
		secretItems.Password = password
	}
	if email != "" {
		secretItems.Email = email
	}
	if notes != "" {
		secretItems.Notes = notes
	}

	// Tags
	if tags != "" {
		rawTags := strings.Split(tags, ",")
		tagList := make([]string, 0, len(rawTags))
		for _, t := range rawTags {
			trimmed := strings.TrimSpace(t)
			if trimmed != "" {
				tagList = append(tagList, trimmed)
			}
		}
		secretItems.Tags = tagList
	}

	if totpSecret != "" {
		secretItems.TOTPSecret = totpSecret
	}

	// Update timestamps
	now := time.Now()
	secretItems.UpdatedAt = now
	vaultItem.UpdatedAt = now

	jsonSecretItems, err := json.Marshal(secretItems)
	if err != nil {
		return err
	}

	cipherText, err := security.EncryptAESGCM(vaultKey, jsonSecretItems, []byte("vault-v1"))
	if err != nil {
		return err
	}

	vaultItem.Data = cipherText

	jsonVaultItem, err := json.Marshal(vaultItem)
	if err != nil {
		return err
	}

	return itemRepo.UpdateData("vault_items", itemName, string(jsonVaultItem))
}
