package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	DBPath              string        `yaml:"db_path"`
	AutoLockTimeout     time.Duration `yaml:"auto_lock_timeout"`
	ClipboardTimeout    time.Duration `yaml:"clipboard_timeout"`
	ShowPasswords       bool          `yaml:"show_passwords"` // Whether to show passwords in plain text
	BackupDir           string        `yaml:"backup_dir"`
	BackupCount         int           `yaml:"backup_count"` // Number of backups to keep
	PasswordLength      int           `yaml:"password_length"`
	UseSpecialChars     bool          `yaml:"use_special_chars"`
	UseNumbers          bool          `yaml:"use_numbers"`
	UseUppercase        bool          `yaml:"use_uppercase"`
	MinPasswordStrength int           `yaml:"min_password_strength"` // 0-4, 0 = no minimum
}

// DefaultConfig returns a new Config with default values
func DefaultConfig() *Config {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}

	defaultDBPath := filepath.Join(home, ".psst", "vault.db")
	defaultBackupDir := filepath.Join(home, ".psst", "backups")

	return &Config{
		DBPath:              defaultDBPath,
		AutoLockTimeout:     15 * time.Minute,
		ClipboardTimeout:    30 * time.Second,
		ShowPasswords:       false,
		BackupDir:           defaultBackupDir,
		BackupCount:         5,
		PasswordLength:      16,
		UseSpecialChars:     true,
		UseNumbers:          true,
		UseUppercase:        true,
		MinPasswordStrength: 2,
	}
}

// LoadConfig loads config from file or creates a default one
func LoadConfig(path string) *Config {
	config := DefaultConfig()

	// Check if config file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Create default config file
		if err := saveConfig(config, path); err != nil {
			fmt.Printf("Error creating default config: %s\n", err)
		}
		return config
	}

	// Read config file
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		return config
	}

	// Parse YAML
	if err := yaml.Unmarshal(data, config); err != nil {
		fmt.Printf("Error parsing config file: %s\n", err)
		return config
	}

	return config
}

// SaveConfig saves the config to file
func saveConfig(config *Config, path string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	// Marshal to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(path, data, 0600)
}
