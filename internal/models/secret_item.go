package models

import "time"

type SecretItem struct {
	ServiceName string   `json:"service_name"`
	Website     string   `json:"website"`
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	Email       string   `json:"email"`
	Notes       string   `json:"notes"`
	Tags        []string `json:"tags"`
	TOTPSecret  string   `json:"totp_secret"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
