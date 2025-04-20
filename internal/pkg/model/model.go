// Package model contains the data model for the vault.
package model

import "time"

// PasswordEntry represents a single password entry in the vault.
type PasswordEntry struct {
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
	LastUsedAt time.Time `json:"last_used_at"`
	Service    string    `json:"service"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	URL        string    `json:"url"`
	Notes      string    `json:"notes"`
	Tags       []string  `json:"tags"`
	ID         int64     `json:"id"`
}

// VaultMetadata represents vault metadata.
type VaultMetadata struct {
	MasterHash string
	CreatedAt  time.Time
	LastAccess time.Time
	Version    string
}
