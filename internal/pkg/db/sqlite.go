package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/CanobbioE/please-safely-store-this/internal/pkg/vault"

	_ "github.com/mattn/go-sqlite3"
)

// Database represents the SQLite database
type Database struct {
	db *sql.DB
}

// NewDatabase creates a new database connection
func NewDatabase(dbPath string) (*Database, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open SQLite database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &Database{db: db}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// Initialize creates the database schema
func (d *Database) Initialize() error {
	// Create password entries table
	_, err := d.db.Exec(`
        CREATE TABLE IF NOT EXISTS password_entries (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            service TEXT NOT NULL,
            username TEXT,
            password TEXT NOT NULL,
            url TEXT,
            notes TEXT,
            created_at TIMESTAMP NOT NULL,
            modified_at TIMESTAMP NOT NULL,
            last_used_at TIMESTAMP
        );
    `)
	if err != nil {
		return fmt.Errorf("failed to create password_entries table: %w", err)
	}

	// Create tags table
	_, err = d.db.Exec(`
        CREATE TABLE IF NOT EXISTS tags (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            entry_id INTEGER,
            tag TEXT NOT NULL,
            FOREIGN KEY (entry_id) REFERENCES password_entries(id) ON DELETE CASCADE
        );
    `)
	if err != nil {
		return fmt.Errorf("failed to create tags table: %w", err)
	}

	// Create vault metadata table
	_, err = d.db.Exec(`
        CREATE TABLE IF NOT EXISTS vault_metadata (
            key TEXT PRIMARY KEY,
            value TEXT NOT NULL
        );
    `)
	if err != nil {
		return fmt.Errorf("failed to create vault_metadata table: %w", err)
	}

	// Create indexes
	_, err = d.db.Exec(`
        CREATE UNIQUE INDEX IF NOT EXISTS idx_password_entries_service ON password_entries(service);
        CREATE INDEX IF NOT EXISTS idx_tags_entry_id ON tags(entry_id);
        CREATE INDEX IF NOT EXISTS idx_tags_tag ON tags(tag);
    `)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}

// SaveVaultMetadata saves vault metadata to the database
func (d *Database) SaveVaultMetadata(v *vault.Vault) error {
	// Insert or update vault metadata
	_, err := d.db.Exec(`
        INSERT OR REPLACE INTO vault_metadata (key, value) VALUES 
        ('master_hash', ?),
        ('created_at', ?),
        ('last_access', ?),
        ('version', ?)
    `, v.MasterHash, v.CreatedAt, v.LastAccess, v.Version)

	if err != nil {
		return fmt.Errorf("failed to save vault metadata: %w", err)
	}

	return nil
}

// SavePasswordEntry saves a password entry to the database
func (d *Database) SavePasswordEntry(entry *vault.PasswordEntry) error {
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert or update password entry
	var result sql.Result
	if entry.ID == 0 {
		// New entry
		result, err = tx.Exec(`
            INSERT INTO password_entries 
            (service, username, password, url, notes, created_at, modified_at, last_used_at)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?)
        `, entry.Service, entry.Username, entry.Password, entry.URL, entry.Notes,
			entry.CreatedAt, entry.ModifiedAt, entry.LastUsedAt)
	} else {
		// Update existing entry
		result, err = tx.Exec(`
            UPDATE password_entries SET
            username = ?, password = ?, url = ?, notes = ?,
            modified_at = ?, last_used_at = ?
            WHERE id = ?
        `, entry.Username, entry.Password, entry.URL, entry.Notes,
			entry.ModifiedAt, entry.LastUsedAt, entry.ID)
	}

	if err != nil {
		return fmt.Errorf("failed to save password entry: %w", err)
	}

	// Get ID for new entry
	if entry.ID == 0 {
		id, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert ID: %w", err)
		}
		entry.ID = id
	}

	// Delete existing tags
	_, err = tx.Exec("DELETE FROM tags WHERE entry_id = ?", entry.ID)
	if err != nil {
		return fmt.Errorf("failed to delete existing tags: %w", err)
	}

	// Insert tags
	for _, tag := range entry.Tags {
		if strings.TrimSpace(tag) == "" {
			continue
		}
		_, err = tx.Exec("INSERT INTO tags (entry_id, tag) VALUES (?, ?)", entry.ID, tag)
		if err != nil {
			return fmt.Errorf("failed to insert tag: %w", err)
		}
	}

	return tx.Commit()
}

// GetPasswordEntry retrieves a password entry by service name
func (d *Database) GetPasswordEntry(service string) (*vault.PasswordEntry, error) {
	var entry vault.PasswordEntry

	// Get password entry
	err := d.db.QueryRow(`
        SELECT id, service, username, password, url, notes, created_at, modified_at, last_used_at
        FROM password_entries
        WHERE service = ?
    `, service).Scan(
		&entry.ID, &entry.Service, &entry.Username, &entry.Password, &entry.URL, &entry.Notes,
		&entry.CreatedAt, &entry.ModifiedAt, &entry.LastUsedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get password entry: %w", err)
	}

	// Get tags
	rows, err := d.db.Query("SELECT tag FROM tags WHERE entry_id = ?", entry.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tags: %w", err)
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, tag)
	}

	entry.Tags = tags

	return &entry, nil
}

// ListPasswordEntries lists all password entries
func (d *Database) ListPasswordEntries() ([]vault.PasswordEntry, error) {
	// Get password entries
	rows, err := d.db.Query(`
        SELECT id, service, username, password, url, notes, created_at, modified_at, last_used_at
        FROM password_entries
        ORDER BY service
    `)
	if err != nil {
		return nil, fmt.Errorf("failed to query password entries: %w", err)
	}
	defer rows.Close()

	var entries []vault.PasswordEntry
	for rows.Next() {
		var entry vault.PasswordEntry
		if err := rows.Scan(
			&entry.ID, &entry.Service, &entry.Username, &entry.Password, &entry.URL, &entry.Notes,
			&entry.CreatedAt, &entry.ModifiedAt, &entry.LastUsedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan password entry: %w", err)
		}
		entries = append(entries, entry)
	}

	// Get tags for each entry
	for i := range entries {
		tagRows, err := d.db.Query("SELECT tag FROM tags WHERE entry_id = ?", entries[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed to query tags: %w", err)
		}

		var tags []string
		for tagRows.Next() {
			var tag string
			if err := tagRows.Scan(&tag); err != nil {
				tagRows.Close()
				return nil, fmt.Errorf("failed to scan tag: %w", err)
			}
			tags = append(tags, tag)
		}
		tagRows.Close()

		entries[i].Tags = tags
	}

	return entries, nil
}

// DeletePasswordEntry deletes a password entry by service name
func (d *Database) DeletePasswordEntry(service string) error {
	// Start transaction
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get entry ID
	var id int64
	err = tx.QueryRow("SELECT id FROM password_entries WHERE service = ?", service).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil // Entry not found, nothing to delete
		}
		return fmt.Errorf("failed to get entry ID: %w", err)
	}

	// Delete tags
	_, err = tx.Exec("DELETE FROM tags WHERE entry_id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete tags: %w", err)
	}

	// Delete entry
	_, err = tx.Exec("DELETE FROM password_entries WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete password entry: %w", err)
	}

	return tx.Commit()
}
