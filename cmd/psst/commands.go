package psst

import (
	"fmt"

	"github.com/spf13/cobra"
)

// TODO: use custom logger instead of fmt

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new password entry",
	Long:  `Add a new password entry to the vault.`,
	Run: func(cmd *cobra.Command, args []string) {
		service, _ := cmd.Flags().GetString("service")
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
			fmt.Println("Password flag not set. Interactive password entry will be implemented in Phase 3.")
			return
		}

		fmt.Printf("Adding password for service: %s, username: %s (Implementation pending)\n", service, username)
		fmt.Printf("Tags: %v\n", tags)

	},
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrieve a password",
	Long:  `Retrieve a password from the vault.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Error: Service name required")
			return
		}
		service := args[0]
		fmt.Printf("Getting password for service: %s (Implementation pending)\n", service)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all password entries",
	Long:  `List all password entries in the vault.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Listing all passwords (Implementation pending)")
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing password",
	Long:  `Update an existing password in the vault.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Error: Service name required")
			return
		}
		service := args[0]
		_, _ = cmd.Flags().GetString("password")
		fmt.Printf("Updating password for service: %s (Implementation pending)\n", service)
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

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a password entry",
	Long:  `Delete a password entry from the vault.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Error: Service name required")
			return
		}
		service := args[0]
		fmt.Printf("Deleting password for service: %s (Implementation pending)\n", service)
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the password vault",
	Long:  `Initialize a new password vault with a master password.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initializing new password vault (Implementation pending)")
	},
}

func init() {
	// Add command flags
	addCmd.Flags().String("service", "", "Service name (required)")
	addCmd.Flags().String("username", "", "Username for the service")
	addCmd.Flags().String("password", "", "Password for the service (if not provided, will prompt)")
	addCmd.Flags().StringSlice("tags", []string{}, "Tags for categorization (comma-separated)")
	addCmd.MarkFlagRequired("service")

	// Update command flags
	updateCmd.Flags().String("password", "", "New password (if not provided, will prompt)")
	updateCmd.Flags().StringSlice("tags", []string{}, "Update tags (comma-separated)")
}
