package secret

import (
	"fmt"
	"log"

	sharedflags "github.com/0xN0nc3Br34k3r/ZeroVaulT/cmd/zerovault/shared"
	"github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/security"
	vaultService "github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/services/vault_item"

	"github.com/spf13/cobra"
)

// CreateSecretItemCommand defines the `zerovault secret create` command.
// It collects user‑provided fields (name, username, password, etc.),
// ensures a master password is available, and delegates creation to the vault service.
var CreateSecretItemCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a new secret item",
	Long: `Create a new secret entry in your ZeroVaulT.

This command allows you to securely store credentials such as usernames,
passwords, emails, TOTP secrets, and additional metadata like tags or notes.
If a master password is not provided via flags, you will be prompted to enter it
securely before the item is created.`,

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

		// Attempt to create the secret item using the vault service.
		// All collected flags are passed directly to the service layer.
		err := vaultService.CreateSecretItem(
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

		// If creation fails, stop execution and show the error.
		if err != nil {
			log.Fatal(err)
			return
		}

		// Success message for the user.
		fmt.Println("Secret item created successfully.")
	},
}

func init() {
	// Register all flags for the `secret create` command.
	// These flags populate the variables used in the Run() function.

	CreateSecretItemCommand.Flags().StringVarP(&ItemName, "name", "n", "", "unique item name")
	CreateSecretItemCommand.Flags().StringVarP(&ServiceName, "service", "s", "", "service name")
	CreateSecretItemCommand.Flags().StringVarP(&Website, "website", "w", "", "website URL")
	CreateSecretItemCommand.Flags().StringVarP(&Username, "username", "u", "", "username")
	CreateSecretItemCommand.Flags().StringVarP(&Password, "password", "p", "", "password")
	CreateSecretItemCommand.Flags().StringVarP(&Email, "email", "e", "", "email address")
	CreateSecretItemCommand.Flags().StringVarP(&Tags, "tags", "t", "", "comma-separated tags")

	// Optional fields without shorthand flags.
	CreateSecretItemCommand.Flags().StringVar(&Notes, "notes", "", "notes")
	CreateSecretItemCommand.Flags().StringVar(&TOTPSecret, "totp", "", "TOTP secret")
}
