package zerovault

import (
	"log"

	sharedflags "github.com/0xN0nc3Br34k3r/ZeroVaulT/cmd/zerovault/shared"

	"github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/security"
	"github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/services/vault"
	"github.com/spf13/cobra"
)

// registerCommand defines the `zerovault register` command.
// This command initializes a new vault by creating and storing
// the master password hash using Argon2id. It must be run once
// before any other commands can be used.
var registerCommand = &cobra.Command{
	Use:   "register",
	Short: "Initialize ZeroVaulT by creating a master password",
	Long: `Register a new master password for ZeroVaulT.

This command must be executed before using any other vault operations.
It securely derives a key from the provided password using Argon2id
and stores the resulting hash in the database. If a vault already
exists, registration will fail to prevent accidental overwrites.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Basic validation: ensure the user actually provided a password.
		if sharedflags.MasterPassword == "" {
			pw, err := security.PromptHidden("Enter master password: ")
			if err != nil {
				log.Fatal(err)
			}
			sharedflags.MasterPassword = pw
		}

		// Attempt to register the vault. If registration fails (e.g., vault already exists),
		// the error is logged and the program exits.
		if err := vault.Register(db, sharedflags.MasterPassword); err != nil {
			log.Fatal(err)
		}

		// If no error occurred, registration succeeded.
		log.Println("Vault successfully initialized.")
	},
}
