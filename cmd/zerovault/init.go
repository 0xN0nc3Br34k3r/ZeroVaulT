package zerovault

import (
	"log"

	"github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/services"
	"github.com/spf13/cobra"
)

// initCommand defines the `zerovault init` command.
// This command prepares the application environment by:
//   - creating the OS‑specific app data directory,
//   - ensuring the database file exists,
//   - initializing required buckets.
//
// It is intended to be run once before using other commands,
// but running it multiple times is safe (idempotent).
var initCommand = &cobra.Command{
	Use:   "init",
	Short: "Initialize ZeroVaulT storage and database",
	Long: `Initialize the ZeroVaulT environment.

This command sets up the application data directory and prepares
the internal database structure. It is safe to run multiple times,
as all operations are idempotent.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Perform environment and database initialization.
		if err := services.Initialization(); err != nil {
			log.Fatal(err)
		}
	},
}
