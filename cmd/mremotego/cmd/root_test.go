package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/jaydenthorup/mremotego/internal/config"
	"github.com/jaydenthorup/mremotego/pkg/models"
)

// setupTestConfig creates a temporary config file for testing
func setupTestConfig(t *testing.T) (string, *config.Manager) {
	t.Helper()

	// Create temp directory
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Create manager and initialize with basic config
	manager := config.NewManager(configPath)
	cfg := models.NewConfig()

	// Add a test folder with connections
	folder := models.NewFolder("Test Folder")

	sshConn := models.NewConnection("Test SSH", models.ProtocolSSH)
	sshConn.Host = "test.example.com"
	sshConn.Port = 22
	sshConn.Username = "testuser"
	sshConn.Description = "Test SSH connection"
	folder.AddChild(sshConn)

	rdpConn := models.NewConnection("Test RDP", models.ProtocolRDP)
	rdpConn.Host = "test-server.local"
	rdpConn.Port = 3389
	rdpConn.Username = "admin"
	folder.AddChild(rdpConn)

	cfg.Connections = append(cfg.Connections, folder)

	if err := manager.Save(); err != nil {
		t.Fatalf("Failed to save test config: %v", err)
	}

	return configPath, manager
}

// captureOutput captures stdout/stderr during command execution
func captureOutput(t *testing.T, fn func()) string {
	t.Helper()

	// Save original stdout
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	// Create pipe to capture output
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	os.Stdout = w
	os.Stderr = w

	// Run function
	fn()

	// Restore stdout/stderr
	w.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	// Read captured output
	var buf bytes.Buffer
	buf.ReadFrom(r)

	return buf.String()
}

func TestRootCommand_Properties(t *testing.T) {
	// Test basic command properties without executing
	if rootCmd.Use != "mremotego" {
		t.Errorf("Expected Use 'mremotego', got %q", rootCmd.Use)
	}
	
	if rootCmd.Short == "" {
		t.Error("Expected Short description to be set")
	}
	
	if rootCmd.Long == "" {
		t.Error("Expected Long description to be set")
	}
}

func TestInitConfig(t *testing.T) {
	// Save original cfgFile
	originalCfgFile := cfgFile
	defer func() { cfgFile = originalCfgFile }()
	
	// Reset cfgFile to test default path resolution
	cfgFile = ""
	
	initConfig()
	
	if cfgFile == "" {
		t.Error("initConfig() did not set cfgFile")
	}
}