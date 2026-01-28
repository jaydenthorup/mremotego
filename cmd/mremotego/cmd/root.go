package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/jaydenthorup/mremotego/internal/config"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "mremotego",
		Short: "A git-compatible remote connection manager",
		Long: `MremoteGO is a Go implementation of mRemoteNG with git-compatible 
configuration files. Store and manage your remote connections (SSH, RDP, VNC, etc.) 
in a human-readable YAML format that works great with version control.`,
	}
)

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/mremotego/config.yaml)")
}

func initConfig() {
	if cfgFile == "" {
		var err error
		cfgFile, err = config.GetDefaultConfigPath()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting default config path: %v\n", err)
			os.Exit(1)
		}
	}
}

// getConfigManager returns a config manager instance
func getConfigManager() (*config.Manager, error) {
	if cfgFile == "" {
		initConfig()
	}

	manager := config.NewManager(cfgFile)
	if err := manager.Load(); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return manager, nil
}
