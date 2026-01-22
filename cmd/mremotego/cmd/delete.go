package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [connection name]",
	Short: "Delete a connection",
	Long:  `Remove a connection from the configuration.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		connectionName := args[0]

		manager, err := getConfigManager()
		if err != nil {
			return err
		}

		// Delete
		if err := manager.DeleteConnection(connectionName); err != nil {
			return fmt.Errorf("failed to delete connection: %w", err)
		}

		// Save
		if err := manager.Save(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("✓ Deleted connection '%s'\n", connectionName)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
