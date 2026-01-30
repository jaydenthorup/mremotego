package cmd

import (
	"strings"
	"testing"

	"github.com/jaydenthorup/mremotego/pkg/models"
)

// resetAddFlags resets all add command flags to their defaults
func resetAddFlags() {
	addName = ""
	addProtocol = ""
	addHost = ""
	addPort = 0
	addUsername = ""
	addPassword = ""
	addDomain = ""
	addDescription = ""
	addFolder = ""
	addTags = []string{}
}

func TestAddCommand_Success(t *testing.T) {
	resetAddFlags()
	configPath, manager := setupTestConfig(t)
	cfgFile = configPath

	// Set flags directly instead of using SetArgs
	addName = "New SSH Server"
	addProtocol = "ssh"
	addHost = "newserver.example.com"
	addPort = 2222
	addUsername = "newuser"
	addDescription = "A new SSH connection"

	// Execute the command function directly
	err := addCmd.RunE(addCmd, []string{})
	if err != nil {
		t.Fatalf("Add command failed: %v", err)
	}

	// Reload config and verify
	if err := manager.Load(); err != nil {
		t.Fatalf("Failed to reload config: %v", err)
	}

	cfg := manager.GetConfig()

	// Should have added one connection - verify it exists
	conn := findConnection(cfg.Connections, "New SSH Server")
	if conn == nil {
		t.Error("Expected to find 'New SSH Server' connection")
	}

	if conn != nil {
		if conn.Host != "newserver.example.com" {
			t.Errorf("Expected host 'newserver.example.com', got %q", conn.Host)
		}
		if conn.Port != 2222 {
			t.Errorf("Expected port 2222, got %d", conn.Port)
		}
	}
}

func TestAddCommand_MissingName(t *testing.T) {
	resetAddFlags()
	configPath, _ := setupTestConfig(t)
	cfgFile = configPath

	// Only set protocol and host, leave name empty
	addProtocol = "ssh"
	addHost = "test.com"

	err := addCmd.RunE(addCmd, []string{})
	if err == nil {
		t.Error("Expected error for missing name, got nil")
		return
	}

	if !strings.Contains(err.Error(), "name is required") {
		t.Errorf("Expected 'name is required' error, got: %v", err)
	}
}

func TestAddCommand_MissingHost(t *testing.T) {
	resetAddFlags()
	configPath, _ := setupTestConfig(t)
	cfgFile = configPath

	// Only set name and protocol, leave host empty
	addName = "Test Connection"
	addProtocol = "ssh"

	err := addCmd.RunE(addCmd, []string{})
	if err == nil {
		t.Error("Expected error for missing host, got nil")
		return
	}

	if !strings.Contains(err.Error(), "host is required") {
		t.Errorf("Expected 'host is required' error, got: %v", err)
	}
}

func TestAddCommand_MissingProtocol(t *testing.T) {
	resetAddFlags()
	configPath, _ := setupTestConfig(t)
	cfgFile = configPath

	// Only set name and host, leave protocol empty
	addName = "Test Connection"
	addHost = "test.com"

	err := addCmd.RunE(addCmd, []string{})
	if err == nil {
		t.Error("Expected error for missing protocol, got nil")
		return
	}

	if !strings.Contains(err.Error(), "protocol is required") {
		t.Errorf("Expected 'protocol is required' error, got: %v", err)
	}
}

func TestAddCommand_WithDefaultPort(t *testing.T) {
	resetAddFlags()
	configPath, manager := setupTestConfig(t)
	cfgFile = configPath

	// Add SSH connection without specifying port
	addName = "SSH Default Port"
	addProtocol = "ssh"
	addHost = "default.example.com"

	err := addCmd.RunE(addCmd, []string{})
	if err != nil {
		t.Fatalf("Add command failed: %v", err)
	}

	// Reload and verify default port was used
	if err := manager.Load(); err != nil {
		t.Fatalf("Failed to reload config: %v", err)
	}

	cfg := manager.GetConfig()
	conn := findConnection(cfg.Connections, "SSH Default Port")

	if conn == nil {
		t.Fatal("Connection not found")
	}

	if conn.Port != 22 {
		t.Errorf("Expected default SSH port 22, got %d", conn.Port)
	}
}

func TestAddCommand_WithCustomPort(t *testing.T) {
	resetAddFlags()
	configPath, manager := setupTestConfig(t)
	cfgFile = configPath

	// Add connection with custom port
	addName = "Custom Port SSH"
	addProtocol = "ssh"
	addHost = "custom.example.com"
	addPort = 2222

	err := addCmd.RunE(addCmd, []string{})
	if err != nil {
		t.Fatalf("Add command failed: %v", err)
	}

	// Reload and verify custom port
	if err := manager.Load(); err != nil {
		t.Fatalf("Failed to reload config: %v", err)
	}

	cfg := manager.GetConfig()
	conn := findConnection(cfg.Connections, "Custom Port SSH")

	if conn == nil {
		t.Fatal("Connection not found")
	}

	if conn.Port != 2222 {
		t.Errorf("Expected custom port 2222, got %d", conn.Port)
	}
}

func TestAddCommand_WithAllFields(t *testing.T) {
	resetAddFlags()
	configPath, manager := setupTestConfig(t)
	cfgFile = configPath

	// Add connection with all fields
	addName = "Complete Connection"
	addProtocol = "rdp"
	addHost = "complete.example.com"
	addPort = 3390
	addUsername = "admin"
	addPassword = "secret"
	addDomain = "EXAMPLE"
	addDescription = "A complete connection with all fields"
	addTags = []string{"production", "windows", "server"}

	err := addCmd.RunE(addCmd, []string{})
	if err != nil {
		t.Fatalf("Add command failed: %v", err)
	}

	// Reload and verify all fields
	if err := manager.Load(); err != nil {
		t.Fatalf("Failed to reload config: %v", err)
	}

	cfg := manager.GetConfig()
	conn := findConnection(cfg.Connections, "Complete Connection")

	if conn == nil {
		t.Fatal("Connection not found")
	}

	// Verify all fields
	if conn.Protocol != "rdp" {
		t.Errorf("Expected protocol 'rdp', got %q", conn.Protocol)
	}
	if conn.Host != "complete.example.com" {
		t.Errorf("Expected host 'complete.example.com', got %q", conn.Host)
	}
	if conn.Port != 3390 {
		t.Errorf("Expected port 3390, got %d", conn.Port)
	}
	if conn.Username != "admin" {
		t.Errorf("Expected username 'admin', got %q", conn.Username)
	}
	if conn.Password != "secret" {
		t.Errorf("Expected password 'secret', got %q", conn.Password)
	}
	if conn.Domain != "EXAMPLE" {
		t.Errorf("Expected domain 'EXAMPLE', got %q", conn.Domain)
	}
	if conn.Description != "A complete connection with all fields" {
		t.Errorf("Expected description, got %q", conn.Description)
	}
	if len(conn.Tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(conn.Tags))
	}
}

func TestAddCommand_Help(t *testing.T) {
	// Verify command properties
	if addCmd.Use != "add" {
		t.Errorf("Expected Use 'add', got %q", addCmd.Use)
	}

	if !strings.Contains(addCmd.Short, "Add") {
		t.Errorf("Expected Short description to contain 'Add', got %q", addCmd.Short)
	}
}

// Helper functions

func countAllConnections(connections []*models.Connection) int {
	count := 0
	for _, conn := range connections {
		if conn.IsFolder() {
			count += countAllConnections(conn.Children)
		} else {
			count++
		}
	}
	return count
}

func findConnection(connections []*models.Connection, name string) *models.Connection {
	for _, conn := range connections {
		if conn.Name == name {
			return conn
		}
		if conn.IsFolder() {
			if found := findConnection(conn.Children, name); found != nil {
				return found
			}
		}
	}
	return nil
}
