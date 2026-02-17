package secret

import "github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/repository"

// db holds a shared database instance for the `secret` command package.
// It is intentionally package‑level so all subcommands can access the same
// initialized repository without passing it through every function.
var db *repository.Database

// InjectDB provides dependency injection for the database layer.
// The root command initializes the database once and injects it here,
// allowing all secret‑related commands to operate on the same connection.
func InjectDB(database *repository.Database) {
	db = database
}
