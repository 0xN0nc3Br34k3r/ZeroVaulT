package vault

import (
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
