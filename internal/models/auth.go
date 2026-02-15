package models

import "time"

type Auth struct {
	PasswordHash string `json:"password_hash"`
	KDFSalt      []byte `json:"kdf_salt"`
	KeyVersion   uint   `json:"key_version"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
