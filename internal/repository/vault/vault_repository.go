package vault

import (
	"fmt"

	"github.com/0xN0nc3Br34k3r/ZeroVaulT/internal/repository"
	"go.etcd.io/bbolt"
)

// VaultRepository provides high‑level operations for managing
// vault‑related data stored inside BoltDB buckets. It acts as
// a domain‑specific layer on top of the generic Database wrapper.
type VaultRepository struct {
	db *repository.Database
}

// NewVaultRepository creates a new VaultRepository instance using
// the provided Database. It does not open or close the database itself.
func NewVaultRepository(db *repository.Database) *VaultRepository {
	return &VaultRepository{db: db}
}

// CreateBucket ensures that a bucket with the given name exists.
// If the bucket does not exist, it will be created. If it already
// exists, the operation is a no‑op. Returns an error if the write
// transaction fails.
func (v *VaultRepository) CreateBucket(bucketName string) error {
	err := v.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})

	return err
}

// InsertData inserts a new key-value pair into the specified bucket.
// The operation fails if the bucket does not exist or if the key
// already exists. Returns an error if the write transaction fails.
func (v *VaultRepository) InsertData(bucketName, key, value string) error {
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
