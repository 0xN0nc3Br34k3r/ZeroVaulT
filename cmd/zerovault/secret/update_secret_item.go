package secret

import (
	"fmt"
	"log"

	sharedflags "github.com/0xN0nc3Br34k3r/ZeroVaulT/cmd/zerovault/shared"
	"github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/security"
	vaultItemService "github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/services/vault_item"
	"github.com/spf13/cobra"
)

// UpdateSecretItemCommand defines the `zerovault secret update` command.
// It modifies an existing secret item stored in the vault, ensuring the
// user provides a valid master password before any update is applied.
//
// This command retrieves the encrypted vault entry associated with the
// provided item name, derives the vault key using the master password,
// decrypts the stored record, applies all user‑supplied field updates,
// re‑encrypts the modified SecretItem, and writes it back to the vault.
//
// If the master password is not supplied via flags, the user will be
// securely prompted to enter it before the update operation begins.
//
// Only the fields explicitly provided via flags will be updated. All
// unspecified fields remain unchanged.
//
// Errors are shown immediately if:
//   - the master password is missing or incorrect
//   - the requested item does not exist
//   - decryption or re‑encryption fails
//   - any internal service error occurs
var UpdateSecretItemCommand = &cobra.Command{
	Use:   "update",
	Short: "Update an existing secret item",
	Long: `Update a stored secret entry in your ZeroVaulT.

This command locates an existing encrypted vault item by its name,
verifies your master password, decrypts the stored data, applies any
field updates you provide, and securely re‑encrypts the modified item
before saving it back to the vault.

Only the fields you explicitly specify will be updated; all others
remain unchanged.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Ensure the master password is available.
		// If not provided via flag, prompt the user securely.
		if sharedflags.MasterPassword == "" {
			pw, err := security.PromptHidden("Enter master password: ")
			if err != nil {
				log.Fatal(err)
			}
			sharedflags.MasterPassword = pw
		}

		// Attempt to read and decrypt the secret item.
		err := vaultItemService.UpdateSecretItem(
			db,
			sharedflags.MasterPassword,
			ItemName,
			ServiceName,
			Website,
			Username,
			Password,
			Email,
			Notes,
			Tags,
			TOTPSecret,
		)

		// If read fails, stop execution and show the error.
		if err != nil {
			log.Fatal(err)
			return
		}

		fmt.Println("Secret item updated successfully.")
	},
}

func init() {
	// Register all flags for the `secret create` command.
	// These flags populate the variables used in the Run() function.

	UpdateSecretItemCommand.Flags().StringVarP(&ItemName, "name", "n", "", "unique item name")
	UpdateSecretItemCommand.Flags().StringVarP(&ServiceName, "service", "s", "", "service name")
	UpdateSecretItemCommand.Flags().StringVarP(&Website, "website", "w", "", "website URL")
	UpdateSecretItemCommand.Flags().StringVarP(&Username, "username", "u", "", "username")
	UpdateSecretItemCommand.Flags().StringVarP(&Password, "password", "p", "", "password")
	UpdateSecretItemCommand.Flags().StringVarP(&Email, "email", "e", "", "email address")
	UpdateSecretItemCommand.Flags().StringVarP(&Tags, "tags", "t", "", "comma-separated tags")

	// Optional fields without shorthand flags.
	UpdateSecretItemCommand.Flags().StringVar(&Notes, "notes", "", "notes")
	UpdateSecretItemCommand.Flags().StringVar(&TOTPSecret, "totp", "", "TOTP secret")

	UpdateSecretItemCommand.MarkFlagRequired("name")
}
