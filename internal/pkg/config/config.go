// Package config handles the application configuration.
package config

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration.
type Config struct {
	DBPath              string        `yaml:"db_path"`
	BackupDir           string        `yaml:"backup_dir"`
	AutoLockTimeout     time.Duration `yaml:"auto_lock_timeout"`
	ClipboardTimeout    time.Duration `yaml:"clipboard_timeout"`
	BackupCount         int           `yaml:"backup_count"`
	PasswordLength      int           `yaml:"password_length"`
	MinPasswordStrength int           `yaml:"min_password_strength"`
	ShowPasswords       bool          `yaml:"show_passwords"`
	UseSpecialChars     bool          `yaml:"use_special_chars"`
	UseNumbers          bool          `yaml:"use_numbers"`
	UseUppercase        bool          `yaml:"use_uppercase"`
}

// DefaultConfig returns a new Config with default values.
func DefaultConfig() *Config {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return &Config{
		DBPath:              filepath.Join(home, ".psst", "vault.db"),
		AutoLockTimeout:     15 * time.Minute,
		ClipboardTimeout:    30 * time.Second,
		ShowPasswords:       false,
		BackupDir:           filepath.Join(home, ".psst", "backups"),
		BackupCount:         5,
		PasswordLength:      16,
		UseSpecialChars:     true,
		UseNumbers:          true,
		UseUppercase:        true,
		MinPasswordStrength: 2,
	}
}

// LoadConfig loads config from file or creates a default one.
func LoadConfig(path string) *Config {
	config := DefaultConfig()

	// Check if config file exists and create it if it doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := SaveConfig(config, path); err != nil {
			log.Printf("Error creating default config: %s\n", err)
		}
		return config
	}

	//nolint:gosec // this is safe
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Error reading config file: %s\n", err)
		return config
	}

	if err = yaml.Unmarshal(data, config); err != nil {
		log.Printf("Error parsing config file: %s\n", err)
		return config
	}

	return config
}

// SaveConfig saves the config to file.
func SaveConfig(config *Config, path string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}

	// Marshal to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(path, data, 0o600)
}
