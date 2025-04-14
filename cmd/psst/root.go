package psst

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/CanobbioE/please-safely-store-this/internal/pkg/config"
)

var cfgFile string
var cfg *config.Config

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "psst",
	Short: "A secure password manager CLI",
	Long: `A secure command-line password manager that stores your passwords
locally in an encrypted database. Your passwords never leave your machine
except when you explicitly export them.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.psst/config.yaml)")
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose output")

	// Add commands
	RootCmd.AddCommand(addCmd)
	RootCmd.AddCommand(getCmd)
	RootCmd.AddCommand(listCmd)
	RootCmd.AddCommand(updateCmd)
	RootCmd.AddCommand(deleteCmd)
	RootCmd.AddCommand(initCmd)
}

func initConfig() {
	// If config file is specified, use it
	if cfgFile != "" {
		// Use config file from the flag
		cfg = config.LoadConfig(cfgFile)
		return
	}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error finding home directory:", err)
		os.Exit(1)
	}

	// Default config location
	configDir := filepath.Join(home, ".psst")
	configPath := filepath.Join(configDir, "config.yaml")

	// Create config directory if it doesn't exist
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err := os.MkdirAll(configDir, 0700)
		if err != nil {
			fmt.Println("Error creating config directory:", err)
			os.Exit(1)
		}
	}

	// Load or create default config
	cfg = config.LoadConfig(configPath)
}
