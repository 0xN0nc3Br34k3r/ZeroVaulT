package zerovault

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/0xN0nc3Br34k3r/ZeroVaulT/cmd/zerovault/secret"
	"github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/repository"
	"github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/services"
	"github.com/spf13/cobra"
)

// db is the global database handle shared across all commands.
// It is lazily initialized in PersistentPreRun.
var db *repository.Database

// Build metadata — these values are replaced at build time using ldflags.
var (
	Version       = "0.1.0"
	Commit        = "dev"
	BuildDate     = "unknown"
	GoVersion     = runtime.Version()
	Platform      = runtime.GOOS + "/" + runtime.GOARCH
	VaultFormat   = "1"               // Increment when vault schema changes
	Argon2Version = "argon2id v1"     // KDF version used for master password hashing
	AESVersion    = "AES-GCM 256-bit" // Encryption algorithm used for vault data
)

// rootCommand is the main entry point for the ZeroVault CLI.
// All subcommands are attached to this command.
var rootCommand = &cobra.Command{
	Use:   "zerovault",
	Short: "ZeroVault is a secure, offline-first password vault",
	Long: `ZeroVault is a command-line password manager built with strong
encryption (Argon2id + AES-GCM) and local storage using bbolt.
It is designed for security, simplicity, and offline reliability.`,
	Version: Version,

	// If the user runs `zerovault` with no subcommand,
	// show a friendly message instead of doing nothing.
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ZeroVault CLI — use --help to see available commands")
	},
}

func init() {
	// Customize how Cobra prints version information.
	// This is shown when the user runs: zerovault --version
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

	initCommand.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		// override root PersistentPreRun so DB is NOT opened for `init`
	}
	rootCommand.AddCommand(initCommand)

	// PersistentPreRun executes before ANY subcommand.
	// We use it to lazily initialize the database exactly once.
	rootCommand.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		// If already initialized (e.g., nested commands), skip.
		if db != nil {
			return
		}

		// Determine the application data directory (platform‑aware).
		dir, err := services.AppDataDir()
		if err != nil {
			log.Fatal(err)
		}

		// Build the full path to the database file.
		dbPath := services.DatabaseFile(dir)

		// Open (or create) the database.
		db, err = repository.NewDatabase(dbPath)
		if err != nil {
			log.Fatal(err)
		}

		// Inject DB into secret package
		secret.InjectDB(db)
	}

	rootCommand.AddCommand(registerCommand)

	rootCommand.AddCommand(SecretCommand)
	SecretCommand.AddCommand(secret.CreateSecretItemCommand)
}

// Execute is called by main.go and starts the CLI.
// It ensures proper exit codes and clean shutdown.
func Execute() {
	// Run the root command and handle any errors.
	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Cleanly close the database after command execution.
	if db != nil {
		db.Close()
	}
}
