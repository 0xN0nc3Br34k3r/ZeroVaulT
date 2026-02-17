package secret

// The following variables hold flag values for the `secret` subcommands.
// Cobra binds CLI flags directly into these package‑level fields during
// command initialization, allowing all secret‑related commands to access
// user‑provided input (name, username, tags, etc.) without passing them
// through function parameters.
//
// These values are populated at runtime when the user executes commands
// such as `zerovault secret create` or `zerovault secret update`.
var (
	ItemName    string // Unique identifier for the secret item
	ServiceName string // Associated service or application name
	Website     string // Optional website URL for the service
	Username    string // Username for the account
	Password    string // Password or credential value
	Email       string // Email associated with the account
	Notes       string // Free‑form notes for additional context
	Tags        string // Comma‑separated tags for categorization
	TOTPSecret  string // Optional TOTP seed for 2FA support
)
