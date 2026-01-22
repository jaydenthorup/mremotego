package secrets

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/1password/onepassword-sdk-go"
)

// Global client instance (lazy initialized)
var (
	globalClient    *onepassword.Client
	clientInitOnce  sync.Once
	clientInitError error
)

// OnePasswordProvider handles retrieving secrets from 1Password SDK
type OnePasswordProvider struct {
	client *onepassword.Client
}

// NewOnePasswordProvider creates a new 1Password provider using the SDK
func NewOnePasswordProvider() *OnePasswordProvider {
	clientInitOnce.Do(func() {
		// Try service account token first
		token := os.Getenv("OP_SERVICE_ACCOUNT_TOKEN")

		var client *onepassword.Client
		var err error

		if token != "" {
			// Use service account token
			client, err = onepassword.NewClient(
				context.Background(),
				onepassword.WithServiceAccountToken(token),
				onepassword.WithIntegrationInfo("MremoteGO", "1.0.0"),
			)
		} else {
			// Use desktop app integration (biometric unlock)
			client, err = onepassword.NewClient(
				context.Background(),
				onepassword.WithIntegrationInfo("MremoteGO", "1.0.0"),
			)
		}

		if err != nil {
			clientInitError = err
			return
		}

		globalClient = client
	})

	return &OnePasswordProvider{
		client: globalClient,
	}
}

// IsEnabled returns whether 1Password SDK is available
func (p *OnePasswordProvider) IsEnabled() bool {
	return p.client != nil && clientInitError == nil
}

// IsReference checks if a string is a 1Password reference (starts with op://)
func (p *OnePasswordProvider) IsReference(value string) bool {
	return strings.HasPrefix(value, "op://")
}

// ResolveSecret retrieves a secret from 1Password using the SDK
// Reference format: op://vault/item/field
// Example: op://Private/MyServer/password
func (p *OnePasswordProvider) ResolveSecret(reference string) (string, error) {
	if !p.IsEnabled() {
		return "", fmt.Errorf("1Password SDK is not available: %v", clientInitError)
	}

	if !p.IsReference(reference) {
		return "", fmt.Errorf("not a 1Password reference: %s", reference)
	}

	// Use the SDK to resolve the secret reference
	secret, err := p.client.Secrets().Resolve(context.Background(), reference)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve secret: %w", err)
	}

	return secret, nil
}

// ResolveIfReference resolves a value if it's a 1Password reference, otherwise returns it as-is
func (p *OnePasswordProvider) ResolveIfReference(value string) string {
	if !p.IsReference(value) {
		return value
	}

	resolved, err := p.ResolveSecret(value)
	if err != nil {
		// If resolution fails, return empty string
		// This prevents exposing the reference format
		fmt.Printf("Warning: Failed to resolve 1Password reference: %v\n", err)
		return ""
	}

	return resolved
}

// CreateItem creates a new Login item in 1Password using the SDK
// Returns the 1Password reference (op://vault/title/password)
func (p *OnePasswordProvider) CreateItem(vault, title, username, password string) (string, error) {
	if !p.IsEnabled() {
		return "", fmt.Errorf("1Password SDK is not available: %v", clientInitError)
	}

	if vault == "" || title == "" {
		return "", fmt.Errorf("vault and title are required")
	}

	// Create a Login item with username and password fields
	item, err := p.client.Items().Create(context.Background(), onepassword.ItemCreateParams{
		VaultID:  vault,
		Title:    title,
		Category: onepassword.ItemCategoryLogin,
		Fields: []onepassword.ItemField{
			{
				ID:        "username",
				Title:     "username",
				FieldType: onepassword.ItemFieldTypeText,
				Value:     username,
			},
			{
				ID:        "password",
				Title:     "password",
				FieldType: onepassword.ItemFieldTypeConcealed,
				Value:     password,
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create 1Password item: %w", err)
	}

	// Return the reference format
	reference := fmt.Sprintf("op://%s/%s/password", vault, item.Title)
	return reference, nil
}

// ListVaults returns a list of available 1Password vaults using the SDK
func (p *OnePasswordProvider) ListVaults() ([]string, error) {
	if !p.IsEnabled() {
		return nil, fmt.Errorf("1Password SDK is not available: %v", clientInitError)
	}

	vaults, err := p.client.Vaults().List(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to list vaults: %w", err)
	}

	vaultNames := make([]string, len(vaults))
	for i, vault := range vaults {
		vaultNames[i] = vault.Title
	}

	return vaultNames, nil
}
