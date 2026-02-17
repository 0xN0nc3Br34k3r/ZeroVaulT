package vaultitem

import (
	"fmt"

	"github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/repository"
	"go.etcd.io/bbolt"
)

// VaultItemRepository provides high‑level operations for managing
// individual vault items stored inside BoltDB buckets. It acts as a
// domain‑specific layer on top of the generic Database wrapper,
// focusing specifically on CRUD operations for secret items.
type VaultItemRepository struct {
	db *repository.Database
}

// NewVaultItemRepository creates a new VaultItemRepository instance
// using the provided Database. It does not open or close the database
// itself; it simply reuses the existing connection managed elsewhere.
func NewVaultItemRepository(db *repository.Database) *VaultItemRepository {
	return &VaultItemRepository{db: db}
}

// InsertData inserts a new key‑value pair into the specified bucket.
// The operation fails if the bucket does not exist or if the key
// already exists. This ensures that secret items cannot be overwritten
// accidentally. Returns an error if the write transaction fails.
func (v *VaultItemRepository) InsertData(bucketName, key, value string) error {
	return v.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %s does not exist", bucketName)
		}

		// Prevent accidental overwrites
		if bucket.Get([]byte(key)) != nil {
			return fmt.Errorf("key %s already exists", key)
		}

		return bucket.Put([]byte(key), []byte(value))
	})
}

// GetData retrieves the value associated with the given key from the
// specified bucket. This method performs a read‑only transaction and
// returns the stored value as a string. If the bucket or key does not
// exist, an error is returned.
func (v *VaultItemRepository) GetData(bucketName, key string) (string, error) {
	var result string

	// Start a read‑only transaction
	err := v.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %s does not exist", bucketName)
		}

		value := bucket.Get([]byte(key))
		if value == nil {
			return fmt.Errorf("key %s does not exist", key)
		}

		result = string(value)
		return nil
	})

	return result, err
}

// GetAllData returns all key‑value pairs from the specified bucket.
// It performs a read‑only transaction and iterates through the bucket
// using a cursor. Returns an error if the bucket does not exist or
// if the read transaction fails.
func (v *VaultItemRepository) GetAllData(bucketName string) (map[string]string, error) {
	result := make(map[string]string)

	err := v.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %s does not exist", bucketName)
		}

		// Create a cursor to iterate over all key/value pairs
		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			result[string(k)] = string(v)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// Delete removes a specific key from the given bucket.
// It opens a write transaction, verifies that the bucket exists,
// and checks whether the target key is present before attempting
// deletion. Returns an error if the bucket is missing, the key
// does not exist, or if the write transaction fails.
func (v *VaultItemRepository) Delete(bucketName, key string) error {
	err := v.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("bucket %s does not exist", bucketName)
		}

		value := b.Get([]byte(key))
		if value == nil {
			return fmt.Errorf("key %s does not exist", key)
		}

		return b.Delete([]byte(key))

	})

	if err != nil {
		return err
	}

	return nil
}

// UpdateData overwrites the value associated with the given key in the
// specified bucket. It opens a write transaction, verifies that the
// bucket exists, and ensures the target key is already present before
// performing the update.
//
// This prevents accidental creation of new records when the intention
// is to modify an existing one. Returns an error if the bucket is
// missing, the key does not exist, or if the write transaction fails.
func (v *VaultItemRepository) UpdateData(bucketName, key, value string) error {
	return v.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("bucket %s does not exist", bucketName)
		}

		// Ensure the key exists before updating
		if b.Get([]byte(key)) == nil {
			return fmt.Errorf("key %s does not exist", key)
		}

		return b.Put([]byte(key), []byte(value))
	})
}
