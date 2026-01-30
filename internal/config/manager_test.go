package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jaydenthorup/mremotego/internal/secrets"
	"github.com/jaydenthorup/mremotego/pkg/models"
)

func TestNewManager(t *testing.T) {
	manager := NewManager("/tmp/test-config.yaml")

	if manager == nil {
		t.Fatal("expected non-nil manager")
	}

	if manager.GetConfigPath() != "/tmp/test-config.yaml" {
		t.Errorf("expected config path /tmp/test-config.yaml, got %s", manager.GetConfigPath())
	}
}

func TestManager_LoadAndSave(t *testing.T) {
	// Create temp directory for test
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-config.yaml")

	manager := NewManager(configPath)

	// Create a test config
	config := models.NewConfig()
	config.Settings.OnePasswordAccount = "TestAccount"
	config.Settings.VaultNames = map[string]string{
		"vault1": "DevOps",
		"vault2": "Employee",
	}
	config.Connections = []*models.Connection{
		models.NewConnection("TestConn", models.ProtocolSSH),
	}

	manager.config = config

	// Save
	err := manager.Save()
	if err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("config file was not created")
	}

	// Load into new manager
	manager2 := NewManager(configPath)
	err = manager2.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	loadedConfig := manager2.GetConfig()

	// Verify settings were loaded
	if loadedConfig.Settings == nil {
		t.Fatal("settings not loaded")
	}
	if loadedConfig.Settings.OnePasswordAccount != "TestAccount" {
		t.Errorf("expected account TestAccount, got %s", loadedConfig.Settings.OnePasswordAccount)
	}
	if len(loadedConfig.Settings.VaultNames) != 2 {
		t.Errorf("expected 2 vault names, got %d", len(loadedConfig.Settings.VaultNames))
	}
	if loadedConfig.Settings.VaultNames["vault1"] != "DevOps" {
		t.Error("vault names not loaded correctly")
	}

	// Verify connections were loaded
	if len(loadedConfig.Connections) != 1 {
		t.Errorf("expected 1 connection, got %d", len(loadedConfig.Connections))
	}
	if loadedConfig.Connections[0].Name != "TestConn" {
		t.Errorf("expected connection name TestConn, got %s", loadedConfig.Connections[0].Name)
	}
}

func TestManager_SaveVaultNameMappings(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-config.yaml")

	manager := NewManager(configPath)
	manager.config = models.NewConfig()

	mappings := map[string]string{
		"vault1": "DevOps",
		"vault2": "Employee",
		"vault3": "Personal",
	}

	err := manager.SaveVaultNameMappings(mappings)
	if err != nil {
		t.Fatalf("failed to save vault mappings: %v", err)
	}

	// Verify mappings were saved
	if len(manager.config.Settings.VaultNames) != 3 {
		t.Errorf("expected 3 vault names, got %d", len(manager.config.Settings.VaultNames))
	}

	// Load and verify persistence
	manager2 := NewManager(configPath)
	err = manager2.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if len(manager2.config.Settings.VaultNames) != 3 {
		t.Errorf("expected 3 vault names after load, got %d", len(manager2.config.Settings.VaultNames))
	}
	if manager2.config.Settings.VaultNames["vault1"] != "DevOps" {
		t.Error("vault mapping not persisted correctly")
	}
}

func TestManager_SetOnePasswordSDKProvider(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-config.yaml")

	manager := NewManager(configPath)
	manager.config = models.NewConfig()
	manager.config.Settings.VaultNames = map[string]string{
		"vault1": "DevOps",
	}

	// Create mock SDK provider
	provider := &secrets.OnePasswordSDKProvider{}
	
	// Set provider
	manager.SetOnePasswordSDKProvider(provider)

	// Verify vault name map was set in provider
	nameMap := provider.GetVaultNameMap()
	if len(nameMap) != 1 {
		t.Errorf("expected 1 vault name in provider, got %d", len(nameMap))
	}
	if nameMap["vault1"] != "DevOps" {
		t.Error("vault name not set in provider")
	}
}

func TestManager_GetConfig(t *testing.T) {
	manager := NewManager("/tmp/test.yaml")
	
	// NewManager doesn't load, but config might be initialized
	// Just verify GetConfig returns whatever is set
	config := models.NewConfig()
	manager.config = config

	// Should return the config
	retrieved := manager.GetConfig()
	if retrieved == nil {
		t.Error("expected non-nil config")
	}
	if retrieved != config {
		t.Error("returned different config instance")
	}
}

func TestManager_AddConnection(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-config.yaml")

	manager := NewManager(configPath)
	manager.config = models.NewConfig()

	conn := models.NewConnection("TestConn", models.ProtocolSSH)
	conn.Host = "example.com"
	conn.Port = 22

	err := manager.AddConnection(conn, "")
	if err != nil {
		t.Fatalf("failed to add connection: %v", err)
	}

	// Verify connection was added
	if len(manager.config.Connections) != 1 {
		t.Errorf("expected 1 connection, got %d", len(manager.config.Connections))
	}
	if manager.config.Connections[0].Name != "TestConn" {
		t.Errorf("expected connection name TestConn, got %s", manager.config.Connections[0].Name)
	}
}

func TestManager_DeleteConnection(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-config.yaml")

	manager := NewManager(configPath)
	manager.config = models.NewConfig()

	conn := models.NewConnection("TestConn", models.ProtocolSSH)
	manager.config.Connections = []*models.Connection{conn}

	err := manager.DeleteConnection(conn.Name)
	if err != nil {
		t.Fatalf("failed to delete connection: %v", err)
	}

	// Verify connection was deleted
	if len(manager.config.Connections) != 0 {
		t.Errorf("expected 0 connections, got %d", len(manager.config.Connections))
	}
}
