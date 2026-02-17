package zerovault

import "github.com/spf13/cobra"

// SecretCommand is the parent command for all secret‑related operations.
// It acts as a logical namespace under `zerovault secret`, grouping
// subcommands such as `create`, `update`, `delete`, and `list`.
// The actual functionality is implemented in subcommands registered
// during initialization in their respective packages.
var SecretCommand = &cobra.Command{
	Use:   "secret",
	Short: "Manage secret items in your ZeroVaulT",
	Long:  "Provides commands for creating, updating, deleting, and listing secret items stored securely in your ZeroVaulT.",
}
