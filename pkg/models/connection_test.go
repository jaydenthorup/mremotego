package models

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()

	if config == nil {
		t.Fatal("expected non-nil config")
	}

	if config.Version != "1.0" {
		t.Errorf("expected version 1.0, got %s", config.Version)
	}

	if config.Settings == nil {
		t.Error("expected Settings to be initialized")
	}

	if config.Connections == nil {
		t.Error("expected Connections to be initialized")
	}

	if len(config.Connections) != 0 {
		t.Errorf("expected 0 connections, got %d", len(config.Connections))
	}
}

func TestSettings_VaultNames(t *testing.T) {
	settings := &Settings{
		OnePasswordAccount: "TestAccount",
		VaultNames: map[string]string{
			"vault1": "DevOps",
			"vault2": "Employee",
		},
	}

	if settings.OnePasswordAccount != "TestAccount" {
		t.Errorf("expected account TestAccount, got %s", settings.OnePasswordAccount)
	}

	if len(settings.VaultNames) != 2 {
		t.Errorf("expected 2 vault names, got %d", len(settings.VaultNames))
	}

	if settings.VaultNames["vault1"] != "DevOps" {
		t.Errorf("expected vault1 -> DevOps, got %s", settings.VaultNames["vault1"])
	}
}

func TestConfig_DeepCopy(t *testing.T) {
	original := &Config{
		Version: "1.0",
		Settings: &Settings{
			OnePasswordAccount: "TestAccount",
			VaultNames: map[string]string{
				"vault1": "DevOps",
				"vault2": "Employee",
			},
		},
		Connections: []*Connection{
			{
				Name:     "Test Connection",
				Type:     NodeTypeConnection,
				Protocol: ProtocolSSH,
				Host:     "example.com",
				Port:     22,
				Username: "user",
				Password: "pass",
			},
		},
	}

	// Create deep copy
	copied := original.DeepCopy()

	// Verify it's not nil
	if copied == nil {
		t.Fatal("expected non-nil copy")
	}

	// Verify version is copied
	if copied.Version != original.Version {
		t.Errorf("expected version %s, got %s", original.Version, copied.Version)
	}

	// Verify settings are copied
	if copied.Settings == nil {
		t.Fatal("expected Settings to be copied")
	}
	if copied.Settings == original.Settings {
		t.Error("expected Settings to be a new instance, not same reference")
	}
	if copied.Settings.OnePasswordAccount != original.Settings.OnePasswordAccount {
		t.Errorf("expected account %s, got %s", original.Settings.OnePasswordAccount, copied.Settings.OnePasswordAccount)
	}

	// Verify vault names are copied
	if len(copied.Settings.VaultNames) != len(original.Settings.VaultNames) {
		t.Errorf("expected %d vault names, got %d", len(original.Settings.VaultNames), len(copied.Settings.VaultNames))
	}
	if copied.Settings.VaultNames["vault1"] != "DevOps" {
		t.Error("vault names not copied correctly")
	}

	// Verify it's a deep copy - modifying copy doesn't affect original
	copied.Settings.VaultNames["vault3"] = "Personal"
	if _, exists := original.Settings.VaultNames["vault3"]; exists {
		t.Error("modifying copy affected original")
	}

	// Verify connections are copied
	if len(copied.Connections) != len(original.Connections) {
		t.Errorf("expected %d connections, got %d", len(original.Connections), len(copied.Connections))
	}
	if copied.Connections[0] == original.Connections[0] {
		t.Error("expected Connections to be deep copied, not same reference")
	}
	if copied.Connections[0].Name != original.Connections[0].Name {
		t.Error("connection data not copied correctly")
	}
}

func TestConfig_DeepCopy_NilSettings(t *testing.T) {
	original := &Config{
		Version:  "1.0",
		Settings: nil,
		Connections: []*Connection{
			{Name: "Test", Type: NodeTypeConnection},
		},
	}

	copied := original.DeepCopy()

	if copied == nil {
		t.Fatal("expected non-nil copy")
	}

	// Settings should remain nil if original was nil
	if copied.Settings != nil {
		t.Error("expected Settings to remain nil")
	}
}

func TestConnection_IsFolder(t *testing.T) {
	tests := []struct {
		name     string
		connType NodeType
		expected bool
	}{
		{
			name:     "Folder type",
			connType: NodeTypeFolder,
			expected: true,
		},
		{
			name:     "Connection type",
			connType: NodeTypeConnection,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := &Connection{
				Type: tt.connType,
			}

			if conn.IsFolder() != tt.expected {
				t.Errorf("expected IsFolder() = %v, got %v", tt.expected, conn.IsFolder())
			}
		})
	}
}

func TestConnection_AddChild(t *testing.T) {
	folder := NewFolder("TestFolder")
	child := NewConnection("TestConnection", ProtocolSSH)

	folder.AddChild(child)

	if len(folder.Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(folder.Children))
	}

	if folder.Children[0] != child {
		t.Error("child not added correctly")
	}

	// Try adding to non-folder (should not add)
	connection := NewConnection("Test", ProtocolSSH)
	connection.AddChild(child)

	if len(connection.Children) != 0 {
		t.Error("child should not be added to non-folder")
	}
}

func TestNewConnection(t *testing.T) {
	conn := NewConnection("TestConn", ProtocolSSH)

	if conn == nil {
		t.Fatal("expected non-nil connection")
	}

	if conn.Name != "TestConn" {
		t.Errorf("expected name TestConn, got %s", conn.Name)
	}

	if conn.Protocol != ProtocolSSH {
		t.Errorf("expected protocol SSH, got %s", conn.Protocol)
	}

	if conn.Type != NodeTypeConnection {
		t.Errorf("expected type connection, got %s", conn.Type)
	}
}

func TestNewFolder(t *testing.T) {
	folder := NewFolder("TestFolder")

	if folder == nil {
		t.Fatal("expected non-nil folder")
	}

	if folder.Name != "TestFolder" {
		t.Errorf("expected name TestFolder, got %s", folder.Name)
	}

	if folder.Type != NodeTypeFolder {
		t.Errorf("expected type folder, got %s", folder.Type)
	}

	if folder.Children == nil {
		t.Error("expected Children to be initialized")
	}

	if len(folder.Children) != 0 {
		t.Errorf("expected 0 children, got %d", len(folder.Children))
	}
}
