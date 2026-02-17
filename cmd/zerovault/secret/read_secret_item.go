package secret

import (
	"fmt"
	"log"

	sharedflags "github.com/0xN0nc3Br34k3r/ZeroVaulT/cmd/zerovault/shared"
	"github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/security"
	vaultItemService "github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/services/vault_item"

	"github.com/spf13/cobra"
)

// ReadSecretItemCommand defines the `zerovault secret read` command.
// It retrieves and decrypts a stored secret item from the vault,
// ensuring the user provides a valid master password before access.
//
// This command loads the encrypted vault entry associated with the
// provided item key, derives the vault key using the master password,
// decrypts the stored ciphertext, and prints the resulting SecretItem
// in a human‑readable JSON format.
//
// If the master password is not supplied via flags, the user will be
// securely prompted to enter it before the read operation proceeds.
//
// Typical usage:
//
//	zerovault secret read --key <itemName>
//
// Errors are shown immediately if:
//   - the master password is missing or incorrect
//   - the requested item does not exist
//   - decryption fails
//   - any internal service error occurs
var ReadSecretItemCommand = &cobra.Command{
	Use:   "read",
	Short: "Read and decrypt a stored secret item",
	Long: `Read a secret entry from your ZeroVaulT.

This command retrieves an encrypted vault item by its key, verifies your
master password, decrypts the stored data, and prints the resulting
SecretItem in a clean, formatted JSON output.

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

		// Attempt to read and decrypt the secret item.
		err := vaultItemService.ReadSecretItem(
			db,
			sharedflags.MasterPassword,
			ItemName,
		)

		// If read fails, stop execution and show the error.
		if err != nil {
			log.Fatal(err)
			return
		}

		// Success message for the user.
		fmt.Println("Secret item read successfully.")
	},
}

func init() {
	// Register the flag used to specify which item to read.
	ReadSecretItemCommand.Flags().StringVarP(
		&ItemName,
		"key",
		"k",
		"",
		"unique item name",
	)

	ReadSecretItemCommand.MarkFlagRequired("key")

}
