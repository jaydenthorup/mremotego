package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/mremotego/internal/config"
	"github.com/yourusername/mremotego/pkg/models"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new configuration file",
	Long:  `Creates a new configuration file with an example connection.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfgFile == "" {
			initConfig()
		}

		manager := config.NewManager(cfgFile)

		// Create a sample config
		cfg := models.NewConfig()

		// Add sample connections
		folder := models.NewFolder("Examples")

		sshConn := models.NewConnection("Example SSH", models.ProtocolSSH)
		sshConn.Host = "example.com"
		sshConn.Port = 22
		sshConn.Username = "user"
		sshConn.Description = "Example SSH connection"
		folder.AddChild(sshConn)

		rdpConn := models.NewConnection("Example RDP", models.ProtocolRDP)
		rdpConn.Host = "windows-server.local"
		rdpConn.Port = 3389
		rdpConn.Username = "Administrator"
		rdpConn.Description = "Example RDP connection"
		folder.AddChild(rdpConn)

		cfg.Connections = append(cfg.Connections, folder)

		// Save to disk
		if err := manager.Save(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Configuration initialized at: %s\n", cfgFile)
		fmt.Println("Edit the file to add your connections or use the 'add' command.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
