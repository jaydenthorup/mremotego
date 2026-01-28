package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/jaydenthorup/mremotego/pkg/models"
)

var (
	editHost        string
	editPort        int
	editUsername    string
	editPassword    string
	editDomain      string
	editDescription string
	editProtocol    string
)

var editCmd = &cobra.Command{
	Use:   "edit [connection name]",
	Short: "Edit an existing connection",
	Long:  `Modify properties of an existing connection.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		connectionName := args[0]

		manager, err := getConfigManager()
		if err != nil {
			return err
		}

		// Create updates object
		updates := &models.Connection{
			Host:        editHost,
			Port:        editPort,
			Username:    editUsername,
			Password:    editPassword,
			Domain:      editDomain,
			Description: editDescription,
			Protocol:    models.Protocol(editProtocol),
		}

		// Update
		if err := manager.UpdateConnection(connectionName, updates); err != nil {
			return fmt.Errorf("failed to update connection: %w", err)
		}

		// Save
		if err := manager.Save(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("✓ Updated connection '%s'\n", connectionName)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(editCmd)

	editCmd.Flags().StringVar(&editHost, "host", "", "New host address")
	editCmd.Flags().IntVar(&editPort, "port", 0, "New port number")
	editCmd.Flags().StringVar(&editUsername, "username", "", "New username")
	editCmd.Flags().StringVar(&editPassword, "password", "", "New password")
	editCmd.Flags().StringVar(&editDomain, "domain", "", "New domain")
	editCmd.Flags().StringVar(&editDescription, "description", "", "New description")
	editCmd.Flags().StringVar(&editProtocol, "protocol", "", "New protocol")
}
