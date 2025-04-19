// Package vault defines operations on the passwords vault.
package vault

import (
	"time"
)

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

// Vault represents the password vault.
type Vault struct {
	LastAccess time.Time        `json:"last_access"`
	CreatedAt  time.Time        `json:"created_at"`
	ModifiedAt time.Time        `json:"modified_at"`
	MasterHash string           `json:"-"`
	Version    string           `json:"version"`
	DBPath     string           `json:"-"`
	Entries    []*PasswordEntry `json:"entries"`
}

// NewVault creates a new empty vault.
func NewVault(dbPath string) *Vault {
	now := time.Now().UTC()
	return &Vault{
		Entries:    make([]*PasswordEntry, 0),
		LastAccess: now,
		CreatedAt:  now,
		ModifiedAt: now,
		Version:    "1.0.0",
		DBPath:     dbPath,
	}
}

// AddEntry adds a new password entry to the vault.
func (v *Vault) AddEntry(entry *PasswordEntry) {
	entry.CreatedAt = time.Now().UTC()
	entry.ModifiedAt = time.Now().UTC()
	v.Entries = append(v.Entries, entry)
	v.ModifiedAt = time.Now().UTC()
}

// GetEntry retrieves a password entry by service name.
func (v *Vault) GetEntry(service string) *PasswordEntry {
	for i, entry := range v.Entries {
		if entry.Service == service {
			return v.Entries[i]
		}
	}
	return nil
}

// UpdateEntry updates an existing password entry.
func (v *Vault) UpdateEntry(service string, newEntry *PasswordEntry) bool {
	for i, entry := range v.Entries {
		if entry.Service == service {
			v.Entries[i] = newEntry
			v.Entries[i].ModifiedAt = time.Now().UTC()
			v.ModifiedAt = time.Now().UTC()
			return true
		}
	}
	return false
}

// DeleteEntry removes a password entry by service name.
func (v *Vault) DeleteEntry(service string) bool {
	for i, entry := range v.Entries {
		if entry.Service == service {
			// Remove the entry by appending everything before and after it
			v.Entries = append(v.Entries[:i], v.Entries[i+1:]...)
			v.ModifiedAt = time.Now().UTC()
			return true
		}
	}
	return false
}

// NewPasswordEntry creates a new entry from the given parameters.
func NewPasswordEntry(service, username, password, url, notes string, tags []string) *PasswordEntry {
	return &PasswordEntry{
		CreatedAt:  time.Now().UTC(),
		ModifiedAt: time.Now().UTC(),
		Service:    service,
		Username:   username,
		Password:   password,
		URL:        url,
		Notes:      notes,
		Tags:       tags,
	}
}
