package shared

// MasterPassword holds the decrypted master password provided by the user.
// It is populated either through a CLI flag or via a secure prompt, and then
// shared across all commands that require authenticated access to the vault.
var MasterPassword string
