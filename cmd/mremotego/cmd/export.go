package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var exportOutput string

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export configuration to a file",
	Long:  `Export the current configuration to a specified file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Read current config file
		data, err := os.ReadFile(cfgFile)
		if err != nil {
			return fmt.Errorf("failed to read config: %w", err)
		}

		// Write to output file
		if err := os.WriteFile(exportOutput, data, 0600); err != nil {
			return fmt.Errorf("failed to write export file: %w", err)
		}

		fmt.Printf("✓ Exported configuration to '%s'\n", exportOutput)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "connections-export.yaml", "Output file path")
}
