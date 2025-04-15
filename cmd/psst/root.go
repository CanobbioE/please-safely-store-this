// Package psst contains all the base commands for the Please Safely Store This tool.
package psst

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/CanobbioE/please-safely-store-this/internal/pkg/config"
)

var cfgFile string

var cfg *config.Config

// RootCmd represents the base command when called without any subcommands.
func RootCmd() *cobra.Command {
	cobra.OnInitialize(initConfig)
	cmd := &cobra.Command{
		Use:   "psst",
		Short: "A secure password manager CLI",
		Long: `A secure command-line password manager that stores your passwords
locally in an encrypted database. Your passwords never leave your machine
except when you explicitly export them.`,
	}
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.psst/config.yaml)")
	cmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose output")
	cmd.AddCommand(AddCmd())
	cmd.AddCommand(GetCmd())
	cmd.AddCommand(ListCmd())
	cmd.AddCommand(UpdateCmd())
	cmd.AddCommand(DeleteCmd())
	cmd.AddCommand(InitCmd())
	return cmd
}

func initConfig() {
	// If config file is specified, use it
	if cfgFile != "" {
		cfg = config.LoadConfig(cfgFile)
		log.Printf("Using config: %+v", cfg) // TODO: remove
		return
	}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Println("Error finding home directory:", err)
		return
	}

	// Default config location
	configDir := filepath.Join(home, ".psst")
	configPath := filepath.Join(configDir, "config.yaml")

	// Create config directory if it doesn't exist
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err := os.MkdirAll(configDir, 0o700)
		if err != nil {
			log.Printf("Error creating config directory: %v", err)
			return
		}
	}

	// Load or create default config
	cfg = config.LoadConfig(configPath)
}
