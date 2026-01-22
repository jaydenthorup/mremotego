package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/mremotego/pkg/models"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all connections",
	Long:  `Display all configured connections in a tree format.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		manager, err := getConfigManager()
		if err != nil {
			return err
		}

		connections := manager.ListConnections()
		if len(connections) == 0 {
			fmt.Println("No connections configured. Use 'mremotego add' to create one.")
			return nil
		}

		// Print connections
		cfg := manager.GetConfig()
		printConnectionTree(cfg.Connections, 0)

		return nil
	},
}

func printConnectionTree(connections []*models.Connection, level int) {
	indent := strings.Repeat("  ", level)

	for _, conn := range connections {
		if conn.IsFolder() {
			fmt.Printf("%s📁 %s\n", indent, conn.Name)
			printConnectionTree(conn.Children, level+1)
		} else {
			icon := getProtocolIcon(conn.Protocol)
			fmt.Printf("%s%s %s (%s://%s", indent, icon, conn.Name, conn.Protocol, conn.Host)
			if conn.Port != 0 {
				fmt.Printf(":%d", conn.Port)
			}
			fmt.Printf(")\n")

			if conn.Description != "" {
				fmt.Printf("%s   └─ %s\n", indent, conn.Description)
			}
		}
	}
}

func getProtocolIcon(protocol models.Protocol) string {
	switch protocol {
	case models.ProtocolSSH:
		return "🔐"
	case models.ProtocolRDP:
		return "🖥️"
	case models.ProtocolVNC:
		return "📺"
	case models.ProtocolHTTP, models.ProtocolHTTPS:
		return "🌐"
	case models.ProtocolTelnet:
		return "📟"
	default:
		return "🔌"
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
}
