package psst

import (
	"log"

	"github.com/spf13/cobra"
)

// TODO: use custom logger instead of log.Print

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

// InitCmd initializes a new password vault with a master password.
func InitCmd() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize the password vault",
		Long:  `Initialize a new password vault with a master password.`,
		Run: func(_ *cobra.Command, _ []string) {
			log.Println("Initializing new password vault (Implementation pending)")
		},
	}

	return initCmd
}
