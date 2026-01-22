package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/mremotego/pkg/models"
)

var (
	addName        string
	addProtocol    string
	addHost        string
	addPort        int
	addUsername    string
	addPassword    string
	addDomain      string
	addDescription string
	addFolder      string
	addTags        []string
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new connection",
	Long:  `Add a new connection to the configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if addName == "" {
			return fmt.Errorf("connection name is required (--name)")
		}
		if addHost == "" {
			return fmt.Errorf("host is required (--host)")
		}
		if addProtocol == "" {
			return fmt.Errorf("protocol is required (--protocol)")
		}

		manager, err := getConfigManager()
		if err != nil {
			return err
		}

		// Create the connection
		conn := models.NewConnection(addName, models.Protocol(addProtocol))
		conn.Host = addHost
		conn.Username = addUsername
		conn.Password = addPassword
		conn.Domain = addDomain
		conn.Description = addDescription
		conn.Tags = addTags

		// Set port or use default
		if addPort != 0 {
			conn.Port = addPort
		} else {
			conn.Port = conn.Protocol.GetDefaultPort()
		}

		// Add to config
		if err := manager.AddConnection(conn, addFolder); err != nil {
			return fmt.Errorf("failed to add connection: %w", err)
		}

		// Save
		if err := manager.Save(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("✓ Added connection '%s' (%s://%s:%d)\n",
			addName, conn.Protocol, conn.Host, conn.Port)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVar(&addName, "name", "", "Connection name (required)")
	addCmd.Flags().StringVar(&addProtocol, "protocol", "", "Protocol: ssh, rdp, vnc, http, https, telnet (required)")
	addCmd.Flags().StringVar(&addHost, "host", "", "Host address or IP (required)")
	addCmd.Flags().IntVar(&addPort, "port", 0, "Port number (default: protocol default)")
	addCmd.Flags().StringVar(&addUsername, "username", "", "Username")
	addCmd.Flags().StringVar(&addPassword, "password", "", "Password (stored in plain text)")
	addCmd.Flags().StringVar(&addDomain, "domain", "", "Domain (for RDP)")
	addCmd.Flags().StringVar(&addDescription, "description", "", "Connection description")
	addCmd.Flags().StringVar(&addFolder, "folder", "", "Folder path (e.g., 'Production/Servers')")
	addCmd.Flags().StringSliceVar(&addTags, "tags", []string{}, "Tags (comma-separated)")
}
