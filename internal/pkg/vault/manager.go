// Package vault handles storing and retrieving passwords from a vault.
// Creates a level of abstraction between the application and the underlying storage.
package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/CanobbioE/please-safely-store-this/internal/pkg/db"
	"github.com/CanobbioE/please-safely-store-this/internal/pkg/model"
)

// Manager represents a vault manager, which can be used to perform CRUD operations on a vault.
// A vault is an abstraction of the underlying database.
//
// A vault manager is created by calling NewManager() and must be initialized with Init before it can be used.
// After initialization, the vault can be immediately used as if it was already unlocked.
// The vault status can be checked by calling IsUnlocked(), the status can be changed by calling Lock() or Unlock().
// Locking the vault manager is useful when the vault manager is no longer needed.
type Manager struct {
	database   *db.Database
	meta       *model.VaultMetadata
	masterKey  []byte
	isUnlocked bool
}

// NewManager creates a new vault manager without initializing it.
// If the vault manager has been already initialized, the Manager can be used after Unlock() has been called.
// If the vault manager has not been initialised, Init() must be called before any other method.
// Please remember to Close() the vault Manager.
func NewManager(dbPath string) (*Manager, error) {
	// TODO: pass the db as a dependency interface.
	repo, err := db.NewDatabase(dbPath)
	if err != nil {
		return nil, fmt.Errorf("error creating database: %w", err)
	}
	return &Manager{
		database: repo,
	}, nil
}

// Close closes and locks (see Lock()) the vault and the underlying database connection.
// The vault can no longer be used after Close has been called.
func (m *Manager) Close() {
	err := m.database.Close()
	if err != nil {
		log.Printf("Error closing database: %s\n", err)
	}
	m.Lock()
}

// Unlock unlocks the vault using the masterPassword.
// If the vault is already unlocked, Unlock returns true and no error.
// The vault can be locked again using Lock.
//
// Unlock returns true if the vault was unlocked successfully, false otherwise.
// If any error occurs, it is returned to the caller.
func (m *Manager) Unlock(masterPassword string) (bool, error) {
	if m.isUnlocked {
		return true, nil
	}
	metadata, err := m.database.GetVaultMetadata()
	if err != nil {
		return false, fmt.Errorf("failed to get vault metadata: %w", err)
	}

	parts := splitHash(metadata.MasterHash)
	if parts == nil {
		return false, errors.New("invalid master hash format")
	}

	// Verify password
	_, key := hashPassword(masterPassword, parts.salt)
	if !verifyPassword(masterPassword, metadata.MasterHash) {
		return false, nil // Password incorrect
	}

	// Unlock vault
	m.meta = metadata
	metadata.LastAccess = time.Now().UTC()
	m.masterKey = key
	m.isUnlocked = true

	// Update last access time
	if err := m.database.SaveVaultMetadata(m.meta); err != nil {
		return true, fmt.Errorf("failed to update last access time: %w", err)
	}

	return true, nil
}

// Lock locks the vault, removing the encryption key from memory.
// The vault can be unlocked again by calling Unlock().
// Locking the vault is useful when the vault is no longer needed.
func (m *Manager) Lock() {
	m.isUnlocked = false
	m.masterKey = nil
	m.meta = nil
}

// IsUnlocked returns true if the vault is unlocked.
// If the vault is locked, the vault can be unlocked by calling Unlock().
// Locking the vault is useful when the vault is no longer needed.
//
// The vault is unlocked by default when the vault is first initialized with Init.
func (m *Manager) IsUnlocked() bool {
	return m.isUnlocked
}

// Init initializes the vault.
// After initialization, the vault can be immediately used as if it was already unlocked.
func (m *Manager) Init(masterPassword string) error {
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}
	hash, key := hashPassword(masterPassword, salt)

	// Initialize the vault
	m.meta = &model.VaultMetadata{
		MasterHash: hash,
		CreatedAt:  time.Now().UTC(),
		LastAccess: time.Now().UTC(),
		Version:    "0.0.1",
	}
	m.masterKey = key
	m.isUnlocked = true

	// Initialize database schema
	if err := m.database.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize database schema: %w", err)
	}

	// Save vault metadata
	if err := m.database.SaveVaultMetadata(m.meta); err != nil {
		return fmt.Errorf("failed to save vault metadata: %w", err)
	}

	return nil
}

// Create adds a new model.PasswordEntry to the vault.
func (m *Manager) Create(entry *model.PasswordEntry) error {
	if !m.isUnlocked {
		return errors.New("vault is locked, please unlock the vault first")
	}
	entry.CreatedAt = time.Now().UTC()
	entry.ModifiedAt = time.Now().UTC()
	entry.LastUsedAt = time.Time{}

	var err error
	entry.Password, err = m.encryptPassword(entry.Password)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	err = m.database.SavePasswordEntry(entry)
	if err != nil {
		return fmt.Errorf("failed to save password entry: %w", err)
	}
	return nil
}

// Read retrieves the model.PasswordEntry associated with the service from the vault.
func (m *Manager) Read(service string) (*model.PasswordEntry, error) {
	if !m.isUnlocked {
		return nil, errors.New("vault is locked, please unlock the vault first")
	}

	entry, err := m.database.GetPasswordEntry(service)
	if err != nil {
		return nil, fmt.Errorf("failed to get password entry: %w", err)
	}

	entry.Password, err = m.decryptPassword(entry.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt password: %w", err)
	}

	return entry, nil
}

// List retrieves all model.PasswordEntry from the vault.
func (*Manager) List() ([]*model.PasswordEntry, error) {
	return nil, errors.New("not implemented")
}

// Update updates a model.PasswordEntry.
func (*Manager) Update(_ *model.PasswordEntry) error {
	return errors.New("not implemented")
}

// Delete removes a model.PasswordEntry from the vault.
func (*Manager) Delete(_ *model.PasswordEntry) error {
	return errors.New("not implemented")
}

// encryptPassword encrypts a password using AES-GCM.
func (m *Manager) encryptPassword(password string) (string, error) {
	block, err := aes.NewCipher(m.masterKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(password), nil)
	return hex.EncodeToString(ciphertext), nil
}

// decryptPassword decrypts a password using AES-GCM.
func (m *Manager) decryptPassword(encryptedPassword string) (string, error) {
	ciphertext, err := hex.DecodeString(encryptedPassword)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(m.masterKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}
	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
