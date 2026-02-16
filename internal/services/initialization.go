package services

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/repository"
	"github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/repository/vault"
)

// AppDataDir determines the correct application data directory
// based on the user's operating system. It creates the directory
// if it does not already exist, using secure permissions (0700).
// This ensures ZeroVaulT always has a private, OS‑appropriate
// storage location for its database and configuration files.
func AppDataDir() (string, error) {
	appName := "ZeroVaulT"
	var base string

	switch runtime.GOOS {
	case "darwin":
		// macOS: ~/Library/Application Support/ZeroVaulT
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, "Library", "Application Support", appName)

	case "linux":
		// Linux: ~/.config/ZeroVaulT
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".config", appName)

	case "windows":
		// Windows: %APPDATA%\ZeroVaulT
		appData := os.Getenv("APPDATA")
		base = filepath.Join(appData, appName)

	default:
		// Fallback for unknown OS: ~/.ZeroVaulT
		home, _ := os.UserHomeDir()
		base = filepath.Join(home, "."+appName)
	}

	// Ensure the directory exists with strict permissions.
	if err := os.MkdirAll(base, 0700); err != nil {
		return "", err
	}

	return base, nil
}

// DatabaseFile returns the absolute path to the ZeroVaulT database
// file inside the application's data directory.
func DatabaseFile(dir string) string {
	return filepath.Join(dir, "ZeroVault.db")
}

// Initialization prepares the application environment on startup.
// It ensures the data directory exists, initializes the database
// file if needed, and guarantees that all required buckets exist.
// This function is intended to run once per program execution.
func Initialization() error {
	// 1. Ensure the application data directory exists.
	dir, err := AppDataDir()
	if err != nil {
		return fmt.Errorf("failed to create app data directory: %w", err)
	}

	// 2. Determine the database file path.
	dbPath := DatabaseFile(dir)

	// 3. Check whether the database already exists.
	_, statErr := os.Stat(dbPath)
	dbExists := !os.IsNotExist(statErr)

	if dbExists {
		fmt.Println("Database already exists:", dbPath)
	} else {
		fmt.Println("Creating new database:", dbPath)
	}

	// 4. Open (or create) the database.
	// The DB is closed immediately after bucket initialization,
	// because commands will open it again when needed.
	db, err := repository.NewDatabase(dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Initialize required buckets for ZeroVaulT.
	v := vault.NewVaultRepository(db)
	buckets := []string{"auth", "dek", "vault_items"}

	for _, name := range buckets {
		// CreateBucket is idempotent: safe to call every startup.
		if err := v.CreateBucket(name); err != nil {
			return fmt.Errorf("failed to create bucket %q: %w", name, err)
		}
	}

	return nil
}
