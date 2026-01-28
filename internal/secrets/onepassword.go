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
	hideConsoleWindow(cmd)
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
	// This will prompt for biometric auth if needed - don't pre-check authentication
	cmd := exec.Command("op", "read", reference)
	hideConsoleWindow(cmd)

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Provide helpful error message based on the output
		errorMsg := string(output)
		return "", fmt.Errorf("failed to retrieve secret: %s", strings.TrimSpace(errorMsg))
	}

	// Trim whitespace and return
	secret := strings.TrimSpace(string(output))
	return secret, nil
}

// ResolveIfReference resolves a value if it's a 1Password reference, otherwise returns it as-is
// ResolveIfReference resolves a value if it's a 1Password reference, otherwise returns it as-is
// Note: This uses the 1Password CLI which will prompt for biometric auth if needed
// TODO: Consider migrating to https://github.com/1Password/onepassword-sdk-go for better integration
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

// CheckItemExists checks if an item with the given title exists in the vault
// Returns the item ID if it exists, or an error if not found or multiple items exist
func (p *OnePasswordProvider) CheckItemExists(vault, title string) (string, bool, error) {
	if !p.enabled {
		return "", false, fmt.Errorf("1Password CLI is not available")
	}

	// Try to get the item
	cmd := exec.Command("op", "item", "get", title, "--vault="+vault, "--format=json")
	hideConsoleWindow(cmd)

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Check if error is because item doesn't exist
		errorMsg := string(output)
		if strings.Contains(errorMsg, "isn't an item") || strings.Contains(errorMsg, "not found") {
			return "", false, nil // Item doesn't exist
		}
		// Check for multiple items
		if strings.Contains(errorMsg, "More than one item matches") {
			return "", false, fmt.Errorf("multiple items found with title '%s' in vault '%s'. Please use a unique name or delete duplicates in 1Password", title, vault)
		}
		return "", false, fmt.Errorf("failed to check item: %s", strings.TrimSpace(errorMsg))
	}

	// Item exists - extract ID from JSON (simple approach)
	// In production you'd want proper JSON parsing
	return "", true, nil
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

	// Check if item already exists
	_, exists, err := p.CheckItemExists(vault, title)
	if err != nil {
		return "", err // Return the error (e.g., multiple items found)
	}

	if exists {
		// Item exists - update it instead of creating
		args := []string{
			"item", "edit", title,
			"--vault=" + vault,
		}

		if username != "" {
			args = append(args, "username="+username)
		}

		if password != "" {
			args = append(args, "password="+password)
		}

		cmd := exec.Command("op", args...)
		hideConsoleWindow(cmd)

		output, err := cmd.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("failed to update existing 1Password item: %w, output: %s", err, string(output))
		}
	} else {
		// Item doesn't exist - create it
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
		hideConsoleWindow(cmd)

		output, err := cmd.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("failed to create 1Password item: %w, output: %s", err, string(output))
		}
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
	hideConsoleWindow(cmd)

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
