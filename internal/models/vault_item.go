package models

import "time"

type VaultItem struct {
	ID   string `json:"id"`
	Data []byte `json:"Data"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
