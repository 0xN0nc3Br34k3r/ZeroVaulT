package repository

import (
	"time"

	"go.etcd.io/bbolt"
)

// Database wraps a BoltDB instance and provides helper methods
// for executing read and write transactions.
type Database struct {
	db *bbolt.DB
}

// NewDatabase opens (or creates) a BoltDB file at the given path.
// A timeout is applied to avoid indefinite locking if another process
// already holds the database file open.
func NewDatabase(path string) (*Database, error) {
	db, err := bbolt.Open(path, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	return &Database{db: db}, nil
}

// Update executes a read-write transaction. The provided function
// receives a writable transaction and should return an error if the
// operation fails.
func (d *Database) Update(fn func(tx *bbolt.Tx) error) error {
	return d.db.Update(fn)
}

// View executes a read-only transaction. The provided function
// receives a read-only transaction and should return an error if the
// operation fails.
func (d *Database) View(fn func(tx *bbolt.Tx) error) error {
	return d.db.View(fn)
}

// Close closes the underlying BoltDB database file.
func (d *Database) Close() error {
	return d.db.Close()
}
