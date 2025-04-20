// Package db handles storing data to sqlite.
package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/CanobbioE/please-safely-store-this/internal/pkg/model"

	_ "github.com/mattn/go-sqlite3" // blank import here is common practice
)

// TODO: refactor to use interfaces and be testable

// Database represents the SQLite database.
type Database struct {
	db *sql.DB
}

// NewDatabase creates a new database connection.
func NewDatabase(dbPath string) (*Database, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o700); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &Database{db: db}, nil
}

// Close closes the database connection.
func (d *Database) Close() error {
	return d.db.Close()
}

// Initialize creates the database schema.
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

// SavePasswordEntry saves a password entry to the database.
func (d *Database) SavePasswordEntry(entry *model.PasswordEntry) error {
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err = tx.Rollback(); err != nil {
			log.Printf("failed to rollback transaction: %v", err)
		}
	}()

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
		entry.ID, err = result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert ID: %w", err)
		}
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

// GetPasswordEntry retrieves a password entry by service name.
func (d *Database) GetPasswordEntry(service string) (*model.PasswordEntry, error) {
	var entry model.PasswordEntry

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
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, tag)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate tags: %w", err)
	}

	entry.Tags = tags

	return &entry, nil
}

// ListPasswordEntries lists all password entries.
func (d *Database) ListPasswordEntries() ([]model.PasswordEntry, error) {
	// Get password entries
	rows, err := d.db.Query(`
        SELECT id, service, username, password, url, notes, created_at, modified_at, last_used_at
        FROM password_entries
        ORDER BY service
    `)
	if err != nil {
		return nil, fmt.Errorf("failed to query password entries: %w", err)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}()

	var entries []model.PasswordEntry
	for rows.Next() {
		var entry model.PasswordEntry
		if err = rows.Scan(
			&entry.ID, &entry.Service, &entry.Username, &entry.Password, &entry.URL, &entry.Notes,
			&entry.CreatedAt, &entry.ModifiedAt, &entry.LastUsedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan password entry: %w", err)
		}
		entries = append(entries, entry)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate password_entries: %w", err)
	}

	// Get tags for each entry
	for i := range entries {
		entries[i].Tags, err = d.getTags(entries[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get tags: %w", err)
		}
	}

	return entries, nil
}

func (d *Database) getTags(entryID int64) ([]string, error) {
	rows, err := d.db.Query("SELECT tag FROM tags WHERE entry_id = ?", entryID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tags: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, tag)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to query tags: %w", rows.Err())
	}

	return tags, nil
}

// DeletePasswordEntry deletes a password entry by service name.
func (d *Database) DeletePasswordEntry(service string) error {
	// Start transaction
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err = tx.Rollback(); err != nil {
			log.Printf("failed to rollback transaction: %v", err)
		}
	}()

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

// SaveVaultMetadata saves vault metadata to the database.
func (d *Database) SaveVaultMetadata(v *model.VaultMetadata) error {
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

// GetVaultMetadata retrieves vault metadata from the database.
func (d *Database) GetVaultMetadata() (*model.VaultMetadata, error) {
	var metadata model.VaultMetadata
	var createdAtStr, lastAccessStr string

	// Get master hash
	err := d.db.QueryRow("SELECT value FROM vault_metadata WHERE key = 'master_hash'").Scan(&metadata.MasterHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get master hash: %w", err)
	}

	// Get created_at
	err = d.db.QueryRow("SELECT value FROM vault_metadata WHERE key = 'created_at'").Scan(&createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get created_at: %w", err)
	}
	metadata.CreatedAt, err = time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}

	// Get last_access
	err = d.db.QueryRow("SELECT value FROM vault_metadata WHERE key = 'last_access'").Scan(&lastAccessStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get last_access: %w", err)
	}
	metadata.LastAccess, err = time.Parse(time.RFC3339, lastAccessStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse last_access: %w", err)
	}

	// Get version
	err = d.db.QueryRow("SELECT value FROM vault_metadata WHERE key = 'version'").Scan(&metadata.Version)
	if err != nil {
		return nil, fmt.Errorf("failed to get version: %w", err)
	}

	return &metadata, nil
}
