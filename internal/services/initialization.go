package services

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/repository"
)

// AppDataDir returns the OS-specific application data directory
// and ensures it exists with secure permissions.
func AppDataDir() (string, error) {
	appName := "ZeroVaulT"
	var base string

	switch runtime.GOOS {
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, "Library", "Application Support", appName)

	case "linux":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".config", appName)

	case "windows":
		appData := os.Getenv("APPDATA")
		base = filepath.Join(appData, appName)

	default:
		home, _ := os.UserHomeDir()
		base = filepath.Join(home, "."+appName)
	}

	if err := os.MkdirAll(base, 0700); err != nil {
		return "", err
	}

	return base, nil
}

// DatabaseFile returns the full path to the database file.
func DatabaseFile(dir string) string {
	return filepath.Join(dir, "ZeroVault.db")
}

// Initialization prepares the application environment and opens the database.
// It returns the app directory path, the database file path, and the DB instance.
func Initialization() (string, string, *repository.Database, error) {
	// 1. Ensure app data directory exists
	dir, err := AppDataDir()
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to create app data directory: %w", err)
	}

	// 2. Build database file path
	dbPath := filepath.Join(dir, "ZeroVault.db")

	// 3. Check if DB file already exists
	_, statErr := os.Stat(dbPath)
	dbExists := !os.IsNotExist(statErr)

	if dbExists {
		fmt.Println("Database already exists:", dbPath)
	} else {
		fmt.Println("Creating new database:", dbPath)
	}

	// 4. Open (or create) the database
	db, err := repository.NewDatabase(dbPath)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to open database: %w", err)
	}

	return dir, dbPath, db, nil
}
