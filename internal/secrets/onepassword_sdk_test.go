package secrets

import (
	"testing"
)

func TestVaultInfo(t *testing.T) {
	tests := []struct {
		name     string
		vaultID  string
		title    string
		expected VaultInfo
	}{
		{
			name:    "Valid vault info",
			vaultID: "abc123",
			title:   "DevOps",
			expected: VaultInfo{
				ID:    "abc123",
				Title: "DevOps",
			},
		},
		{
			name:    "Encrypted title",
			vaultID: "xyz789",
			title:   "[Encrypted]",
			expected: VaultInfo{
				ID:    "xyz789",
				Title: "[Encrypted]",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vault := VaultInfo{
				ID:    tt.vaultID,
				Title: tt.title,
			}

			if vault.ID != tt.expected.ID {
				t.Errorf("expected ID %s, got %s", tt.expected.ID, vault.ID)
			}
			if vault.Title != tt.expected.Title {
				t.Errorf("expected Title %s, got %s", tt.expected.Title, vault.Title)
			}
		})
	}
}

func TestOnePasswordSDKProvider_SetVaultNameMap(t *testing.T) {
	provider := &OnePasswordSDKProvider{
		vaultNameMap: make(map[string]string),
	}

	mappings := map[string]string{
		"vault1": "DevOps",
		"vault2": "Employee",
		"vault3": "Personal",
	}

	provider.SetVaultNameMap(mappings)

	if len(provider.vaultNameMap) != 3 {
		t.Errorf("expected 3 mappings, got %d", len(provider.vaultNameMap))
	}

	for id, name := range mappings {
		if provider.vaultNameMap[id] != name {
			t.Errorf("expected mapping %s -> %s, got %s", id, name, provider.vaultNameMap[id])
		}
	}
}

func TestOnePasswordSDKProvider_GetVaults(t *testing.T) {
	tests := []struct {
		name         string
		enabled      bool
		vaults       []VaultInfo
		vaultNameMap map[string]string
		expected     []VaultInfo
	}{
		{
			name:    "Not enabled",
			enabled: false,
			vaults: []VaultInfo{
				{ID: "vault1", Title: "[Encrypted]"},
			},
			vaultNameMap: nil,
			expected:     []VaultInfo{},
		},
		{
			name:    "No name mappings",
			enabled: true,
			vaults: []VaultInfo{
				{ID: "vault1", Title: "[Encrypted]"},
				{ID: "vault2", Title: "[Encrypted]"},
			},
			vaultNameMap: nil,
			expected: []VaultInfo{
				{ID: "vault1", Title: "[Encrypted]"},
				{ID: "vault2", Title: "[Encrypted]"},
			},
		},
		{
			name:    "With name mappings",
			enabled: true,
			vaults: []VaultInfo{
				{ID: "vault1", Title: "[Encrypted]"},
				{ID: "vault2", Title: "[Encrypted]"},
			},
			vaultNameMap: map[string]string{
				"vault1": "DevOps",
				"vault2": "Employee",
			},
			expected: []VaultInfo{
				{ID: "vault1", Title: "DevOps"},
				{ID: "vault2", Title: "Employee"},
			},
		},
		{
			name:    "Partial name mappings",
			enabled: true,
			vaults: []VaultInfo{
				{ID: "vault1", Title: "[Encrypted]"},
				{ID: "vault2", Title: "[Encrypted]"},
				{ID: "vault3", Title: "[Encrypted]"},
			},
			vaultNameMap: map[string]string{
				"vault1": "DevOps",
			},
			expected: []VaultInfo{
				{ID: "vault1", Title: "DevOps"},
				{ID: "vault2", Title: "[Encrypted]"},
				{ID: "vault3", Title: "[Encrypted]"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &OnePasswordSDKProvider{
				enabled:      tt.enabled,
				vaults:       tt.vaults,
				vaultNameMap: tt.vaultNameMap,
			}

			result := provider.GetVaults()

			if len(result) != len(tt.expected) {
				t.Fatalf("expected %d vaults, got %d", len(tt.expected), len(result))
			}

			for i, expected := range tt.expected {
				if result[i].ID != expected.ID {
					t.Errorf("vault %d: expected ID %s, got %s", i, expected.ID, result[i].ID)
				}
				if result[i].Title != expected.Title {
					t.Errorf("vault %d: expected Title %s, got %s", i, expected.Title, result[i].Title)
				}
			}
		})
	}
}

func TestOnePasswordSDKProvider_GetVaultNameMap(t *testing.T) {
	provider := &OnePasswordSDKProvider{
		vaultNameMap: map[string]string{
			"vault1": "DevOps",
			"vault2": "Employee",
		},
	}

	result := provider.GetVaultNameMap()

	if len(result) != 2 {
		t.Errorf("expected 2 mappings, got %d", len(result))
	}

	if result["vault1"] != "DevOps" {
		t.Errorf("expected vault1 -> DevOps, got %s", result["vault1"])
	}
	if result["vault2"] != "Employee" {
		t.Errorf("expected vault2 -> Employee, got %s", result["vault2"])
	}
}

func TestOnePasswordSDKProvider_IsEnabled(t *testing.T) {
	tests := []struct {
		name     string
		enabled  bool
		expected bool
	}{
		{
			name:     "Enabled",
			enabled:  true,
			expected: true,
		},
		{
			name:     "Disabled",
			enabled:  false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &OnePasswordSDKProvider{
				enabled: tt.enabled,
			}

			if provider.IsEnabled() != tt.expected {
				t.Errorf("expected IsEnabled() = %v, got %v", tt.expected, provider.IsEnabled())
			}
		})
	}
}
