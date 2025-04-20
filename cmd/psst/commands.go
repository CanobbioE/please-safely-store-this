package psst

import (
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/CanobbioE/please-safely-store-this/internal/pkg/db"
	"github.com/CanobbioE/please-safely-store-this/internal/pkg/vault"
)

// TODO: use custom logger instead of log.Print

// minPrincipalPasswordLength is the minimum length of the principal password.
const minPrincipalPasswordLength = 8

// AddCmd adds a new password entry to the vault.
func AddCmd() *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new password entry",
		Long:  `Add a new password entry to the vault.`,
		Run: func(cmd *cobra.Command, _ []string) {
			service, _ := cmd.Flags().GetString("service")
			if service == "" {
				log.Println("Error: Service name required")
				return
			}
			username, _ := cmd.Flags().GetString("username")
			password, _ := cmd.Flags().GetString("password")
			tags, _ := cmd.Flags().GetStringSlice("tags")

			// TODO:
			// connect to db
			// check if service already exists
			// If password flag not set, prompt for password
			// Check password strength if configured
			// Create new password entry
			// Save to database
			// If configured to show tags, display them

			if password == "" {
				log.Println("Password flag not set. Interactive password entry will be implemented in Phase 3.")
				return
			}

			log.Printf("Adding password for service: %s, username: %s (Implementation pending)\n", service, username)
			log.Printf("Tags: %v\n", tags)
		},
	}
	addCmd.Flags().String("service", "", "Service name (required)")
	addCmd.Flags().String("username", "", "Username for the service")
	addCmd.Flags().String("password", "", "Password for the service (if not provided, will prompt)")
	addCmd.Flags().StringSlice("tags", []string{}, "Tags for categorization (comma-separated)")
	err := addCmd.MarkFlagRequired("service")
	if err != nil {
		log.Println(err)
	}
	return addCmd
}

// GetCmd retrieves a password from the vault.
func GetCmd() *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Retrieve a password",
		Long:  `Retrieve a password from the vault.`,
		Run: func(cmd *cobra.Command, _ []string) {
			service, _ := cmd.Flags().GetString("service")
			if service == "" {
				log.Println("Error: Service name required")
				return
			}
			log.Printf("Getting password for service: %s (Implementation pending)\n", service)
		},
	}
	getCmd.Flags().String("service", "", "Service name (required)")
	err := getCmd.MarkFlagRequired("service")
	if err != nil {
		log.Println(err)
	}
	return getCmd
}

// ListCmd lists all password entries in the vault.
func ListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all password entries",
		Long:  `List all password entries in the vault.`,
		Run: func(_ *cobra.Command, _ []string) {
			log.Println("Listing all passwords (Implementation pending)")
		},
	}

	return listCmd
}

// UpdateCmd updates an existing password in the vault.
func UpdateCmd() *cobra.Command {
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update an existing password",
		Long:  `Update an existing password in the vault.`,
		Run: func(cmd *cobra.Command, _ []string) {
			service, _ := cmd.Flags().GetString("service")
			if service == "" {
				log.Println("Error: Service name required")
				return
			}
			password, _ := cmd.Flags().GetString("password")
			if password == "" {
				log.Println("Password flag not set. Interactive password entry will be implemented in Phase 3.")
				return
			}

			log.Printf("Updating password for service: %s (Implementation pending)\n", service)
			// TODO:
			// connect to db
			// check if service exists
			// If password flag not set, prompt for password
			// Check password strength if configured
			// Create new password entry
			// Save to database
			// If configured to show tags, display them
		},
	}
	updateCmd.Flags().String("service", "", "Service name (required)")
	updateCmd.Flags().String("password", "", "New password (if not provided, will prompt)")
	updateCmd.Flags().StringSlice("tags", []string{}, "Update tags (comma-separated)")
	err := updateCmd.MarkFlagRequired("service")
	if err != nil {
		log.Println(err)
	}
	return updateCmd
}

// DeleteCmd deletes a password entry from the vault.
func DeleteCmd() *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a password entry",
		Long:  `Delete a password entry from the vault.`,
		Run: func(cmd *cobra.Command, _ []string) {
			service, _ := cmd.Flags().GetString("service")
			if service == "" {
				log.Println("Error: Service name required")
				return
			}
			log.Printf("Deleting password for service: %s (Implementation pending)\n", service)
		},
	}
	deleteCmd.Flags().String("service", "", "Service name (required)")
	err := deleteCmd.MarkFlagRequired("service")
	if err != nil {
		log.Println(err)
	}
	return deleteCmd
}

// InitCmd initializes a new password vault with a principal password.
func InitCmd() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize the password vault",
		Long:  `Initialize a new password vault with a principal password.`,
		Run: func(_ *cobra.Command, _ []string) {
			initVaultManager()
			defer vaultManager.Close()
			readPwd := func(prompt string) string {
				log.Print(prompt)
				p, err := term.ReadPassword(syscall.Stdin)
				if err != nil {
					log.Printf("Error reading password: %s\n", err)
					return ""
				}
				log.Println()
				return string(p)
			}

			// If the vault already exists, prompt for confirmation to overwrite it.
			_, err := os.Stat(cfg.DBPath)
			if err != nil && !os.IsNotExist(err) {
				log.Printf("Error checking vault file: %s\n", err)
				return
			}
			if !os.IsNotExist(err) {
				log.Printf("Vault already exists at %s. Do you want to overwrite it? [y/N] ", cfg.DBPath)
				var confirm string
				_, err = fmt.Scanln(&confirm)
				if err != nil {
					log.Printf("Error reading confirmation: %s\n", err)
					return
				}
				if !strings.EqualFold(confirm, "y") {
					log.Println("Initialization cancelled.")
					return
				}
				log.Println("Removing existing vault...")
				// Remove existing vault
				if err = os.Remove(cfg.DBPath); err != nil {
					log.Printf("Error removing existing vault: %s\n", err)
					return
				}
			}

			// Prompt for the principal password and confirm it
			// Retry until the password is at least 8 characters long and matches the confirmation
		promptPwd:
			password := readPwd("Enter principal password: ")
			if len(password) < minPrincipalPasswordLength {
				log.Printf("Password must be at least %d characters long!\n", minPrincipalPasswordLength)
				goto promptPwd
			}

			confirmPassword := readPwd("Confirm principal password: ")
			if password != confirmPassword {
				log.Println("Passwords do not match! Please try again.")
				goto promptPwd
			}

			v, err := db.NewDatabase(cfg.DBPath)
			if err != nil {
				log.Printf("Error creating db connection: %v\n", err)
				return
			}
			vaultManager = vault.NewManager(v)

			log.Println("Initializing vault...")
			err = vaultManager.Init(password)
			if err != nil {
				log.Printf("Error initializing vault: %s\n", err)
				return
			}
			log.Println("Vault initialized successfully with the given principal password.")
			log.Println("Please make sure to store the principal password somewhere safe.")
			log.Println("You can now add passwords to the vault using the 'add' command.")
		},
	}

	return initCmd
}

func initVaultManager() {
	if vaultManager != nil {
		return
	}

	v, err := db.NewDatabase(cfg.DBPath)
	if err != nil {
		log.Printf("Error creating DB connection: %v\n", err)
		return
	}

	vaultManager = vault.NewManager(v)
}
