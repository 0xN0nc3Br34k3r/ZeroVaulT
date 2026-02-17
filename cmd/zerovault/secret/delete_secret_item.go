package secret

import (
	"fmt"
	"log"

	sharedflags "github.com/0xN0nc3Br34k3r/ZeroVaulT/cmd/zerovault/shared"
	"github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/security"
	vaultItemService "github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/services/vault_item"
	"github.com/spf13/cobra"
)

// DeleteSecretItemCommand defines the `zerovault secret delete` command.
// It removes a stored secret item from the vault after verifying the
// user's master password.
//
// This command loads the encrypted vault entry associated with the
// provided item key, derives the vault key using the master password,
// validates access, and deletes the corresponding record from the vault.
//
// If the master password is not supplied via flags, the user will be
// securely prompted to enter it before the delete operation proceeds.
//
// Typical usage:
//
//	zerovault secret delete --key <itemName>
//
// Errors are shown immediately if:
//   - the master password is missing or incorrect
//   - the requested item does not exist
//   - deletion fails
//   - any internal service error occurs
var DeleteSecretItemCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete a stored secret item from the vault",
	Long: `Delete a secret entry from your ZeroVaulT.

This command verifies your master password, locates the encrypted vault
item by its key, and permanently removes it from the vault.

If no master password is provided via flags, you will be prompted to
enter it securely before the delete operation begins.`,
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

		// Attempt to delete the secret item.
		err := vaultItemService.Delete(
			db,
			sharedflags.MasterPassword,
			ItemName,
		)

		// If deletion fails, stop execution and show the error.
		if err != nil {
			log.Fatal(err)
			return
		}

		// Success message for the user.
		fmt.Printf("Item %s successfully deleted.\n", ItemName)
	},
}

func init() {
	DeleteSecretItemCommand.Flags().StringVarP(
		&ItemName,
		"key",
		"k",
		"",
		"unique item name",
	)

	DeleteSecretItemCommand.MarkFlagRequired("key")

}
