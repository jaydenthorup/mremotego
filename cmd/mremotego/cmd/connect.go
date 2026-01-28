package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/jaydenthorup/mremotego/internal/launcher"
)

var connectCmd = &cobra.Command{
	Use:   "connect [connection name]",
	Short: "Connect to a configured host",
	Long:  `Launch a connection using the configured protocol handler.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		connectionName := args[0]

		manager, err := getConfigManager()
		if err != nil {
			return err
		}

		// Find the connection
		conn, err := manager.FindConnection(connectionName)
		if err != nil {
			return fmt.Errorf("connection not found: %w", err)
		}

		// Launch
		l := launcher.NewLauncher()
		if err := l.Launch(conn); err != nil {
			return fmt.Errorf("failed to launch connection: %w", err)
		}

		fmt.Printf("✓ Launched connection to '%s' (%s://%s:%d)\n",
			conn.Name, conn.Protocol, conn.Host, conn.Port)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
