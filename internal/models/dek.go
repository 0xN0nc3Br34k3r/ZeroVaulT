package models

import "time"

type DEK struct {
	EncryptedKey []byte `json:"encrypted_key"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
