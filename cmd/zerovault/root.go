package zerovault

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	Version       = "0.1.0"
	Commit        = "dev"
	BuildDate     = "unknown"
	GoVersion     = runtime.Version()
	Platform      = runtime.GOOS + "/" + runtime.GOARCH
	VaultFormat   = "1" // increment if vault schema changes
	Argon2Version = "argon2id v1"
	AESVersion    = "AES-GCM 256-bit"
)

var rootCommand = &cobra.Command{
	Use:   "zerovault",
	Short: "ZeroVault is a secure, offline-first password vault",
	Long: `ZeroVault is a command-line password manager built with strong
encryption (Argon2id + AES-GCM) and local storage using bbolt.
It is designed for security, simplicity, and offline reliability.`,
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		// Default behavior when no subcommand is provided
		fmt.Println("ZeroVault CLI — use --help to see available commands")
	},
}

func init() {
	// Customize how --version prints
	rootCommand.SetVersionTemplate(
		`ZeroVault    {{.Version}}
Commit:      ` + Commit + `
Built:       ` + BuildDate + `
Go:          ` + GoVersion + `
Platform:    ` + Platform + `
VaultFormat: ` + VaultFormat + `
Crypto:      ` + Argon2Version + `, ` + AESVersion + `
`,
	)

	rootCommand.AddCommand(initCommand)
}

// Execute is called by main.go
func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
