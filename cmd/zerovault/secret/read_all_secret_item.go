package secret

import (
	"fmt"
	"log"

	sharedflags "github.com/0xN0nc3Br34k3r/ZeroVaulT/cmd/zerovault/shared"
	"github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/security"
	vaultItemService "github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/services/vault_item"

	"github.com/spf13/cobra"
)

// ReadAllSecretItemCommand defines the `zerovault secret readall` command.
// It retrieves and decrypts *all* stored secret items from the vault,
// ensuring the user provides a valid master password before access.
//
// This command loads every encrypted VaultItem stored in the "vault_items"
// bucket, derives the vault key using the master password, decrypts each
// entry, and prints all SecretItem objects in formatted JSON.
//
// If the master password is not supplied via flags, the user will be
// securely prompted to enter it before the read operation proceeds.
//
// Typical usage:
//
//	zerovault secret readall
//
// Errors are shown immediately if:
//   - the master password is missing or incorrect
//   - the vault bucket cannot be read
//   - any internal service error occurs
var ReadAllSecretItemCommand = &cobra.Command{
	Use:   "readall",
	Short: "Read and decrypt all stored secret items",
	Long: `Read all secret entries from your ZeroVaulT.

This command retrieves every encrypted vault item, verifies your master
password, decrypts each stored record, and prints all SecretItem objects
in a clean, formatted JSON output.

If no master password is provided via flags, you will be prompted to
enter it securely before the read operation begins.`,
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

		// Attempt to read and decrypt all secret items.
		err := vaultItemService.ReadAllSecretItem(
			db,
			sharedflags.MasterPassword,
		)

		// If read fails, stop execution and show the error.
		if err != nil {
			log.Fatal(err)
			return
		}

		// Success message for the user.
		fmt.Println("Secret items read successfully.")
	},
}
