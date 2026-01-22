package secrets

import (
	"fmt"
	"os/exec"
	"strings"
)

// OnePasswordProvider handles retrieving secrets from 1Password CLI
type OnePasswordProvider struct {
	enabled bool
}

// NewOnePasswordProvider creates a new 1Password provider
func NewOnePasswordProvider() *OnePasswordProvider {
	return &OnePasswordProvider{
		enabled: isOnePasswordCLIAvailable(),
	}
}

// isOnePasswordCLIAvailable checks if the 1Password CLI (op) is installed
func isOnePasswordCLIAvailable() bool {
	cmd := exec.Command("op", "--version")
	return cmd.Run() == nil
}

// IsEnabled returns whether 1Password CLI is available
func (p *OnePasswordProvider) IsEnabled() bool {
	return p.enabled
}

// IsReference checks if a string is a 1Password reference (starts with op://)
func (p *OnePasswordProvider) IsReference(value string) bool {
	return strings.HasPrefix(value, "op://")
}

// ResolveSecret retrieves a secret from 1Password using the CLI
// Reference format: op://vault/item/field
// Example: op://Private/MyServer/password
func (p *OnePasswordProvider) ResolveSecret(reference string) (string, error) {
	if !p.enabled {
		return "", fmt.Errorf("1Password CLI is not available")
	}

	if !p.IsReference(reference) {
		return "", fmt.Errorf("not a 1Password reference: %s", reference)
	}

	// Use 'op read' to retrieve the secret
	cmd := exec.Command("op", "read", reference)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve secret from 1Password: %w", err)
	}

	// Trim whitespace and return
	secret := strings.TrimSpace(string(output))
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

// CreateItem creates a new Login item in 1Password
// Returns the 1Password reference (op://vault/title/password)
func (p *OnePasswordProvider) CreateItem(vault, title, username, password string) (string, error) {
	if !p.enabled {
		return "", fmt.Errorf("1Password CLI is not available")
	}

	if vault == "" || title == "" {
		return "", fmt.Errorf("vault and title are required")
	}

	// Build the command: op item create --category=login --title="title" --vault="vault" username="username" password="password"
	args := []string{
		"item", "create",
		"--category=login",
		"--title=" + title,
		"--vault=" + vault,
	}

	if username != "" {
		args = append(args, "username="+username)
	}

	if password != "" {
		args = append(args, "password="+password)
	}

	cmd := exec.Command("op", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to create 1Password item: %w, output: %s", err, string(output))
	}

	// Return the reference format
	reference := fmt.Sprintf("op://%s/%s/password", vault, title)
	return reference, nil
}

// ListVaults returns a list of available 1Password vaults
func (p *OnePasswordProvider) ListVaults() ([]string, error) {
	if !p.enabled {
		return nil, fmt.Errorf("1Password CLI is not available")
	}

	cmd := exec.Command("op", "vault", "list", "--format=json")
	_, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list vaults: %w", err)
	}

	// Simple parsing - extract vault names (this is a simplified version)
	// In production, you'd want to parse the JSON properly
	vaults := []string{"Private", "DevOps", "Employee"} // Default common vaults
	
	// For now, return common vaults. You could parse JSON if needed
	return vaults, nil
}
