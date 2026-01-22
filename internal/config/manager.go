package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/yourusername/mremotego/internal/launcher"
	"github.com/yourusername/mremotego/internal/secrets"
	"github.com/yourusername/mremotego/pkg/models"
	"gopkg.in/yaml.v3"
)

// Manager handles configuration file operations
type Manager struct {
	configPath          string
	config              *models.Config
	onePasswordProvider *secrets.OnePasswordProvider
}

// NewManager creates a new configuration manager
func NewManager(configPath string) *Manager {
	return &Manager{
		configPath:          configPath,
		onePasswordProvider: secrets.NewOnePasswordProvider(),
	}
}

// GetDefaultConfigPath returns the default configuration file path
func GetDefaultConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	// Use platform-specific config directory
	var configDir string
	if os.Getenv("APPDATA") != "" {
		// Windows
		configDir = filepath.Join(os.Getenv("APPDATA"), "mremotego")
	} else {
		// Linux/Mac
		configDir = filepath.Join(homeDir, ".config", "mremotego")
	}

	return filepath.Join(configDir, "config.yaml"), nil
}

// Load loads the configuration from disk
func (m *Manager) Load() error {
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create a new config if file doesn't exist
			m.config = models.NewConfig()
			return nil
		}
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var config models.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// Decrypt passwords after loading (if on Windows with DPAPI support)
	m.decryptPasswords(&config)

	m.config = &config
	return nil
}

// Save writes the configuration to disk
func (m *Manager) Save() error {
	if m.config == nil {
		m.config = models.NewConfig()
	}

	// Ensure directory exists
	dir := filepath.Dir(m.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create a copy of config for encryption (don't modify original)
	configCopy := *m.config
	m.encryptPasswords(&configCopy)

	data, err := yaml.Marshal(&configCopy)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(m.configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetConfig returns the current configuration
func (m *Manager) GetConfig() *models.Config {
	if m.config == nil {
		m.config = models.NewConfig()
	}
	return m.config
}

// AddConnection adds a new connection to the config
func (m *Manager) AddConnection(conn *models.Connection, folderPath string) error {
	if m.config == nil {
		m.config = models.NewConfig()
	}

	conn.Created = time.Now().Format(time.RFC3339)
	conn.Modified = conn.Created

	if folderPath == "" {
		// Add to root
		m.config.Connections = append(m.config.Connections, conn)
		return nil
	}

	// Find or create folder path
	folder, err := m.findOrCreateFolder(folderPath)
	if err != nil {
		return err
	}

	folder.AddChild(conn)
	return nil
}

// FindConnection finds a connection by name (searches recursively)
func (m *Manager) FindConnection(name string) (*models.Connection, error) {
	if m.config == nil {
		return nil, fmt.Errorf("config not loaded")
	}

	return m.findConnectionRecursive(name, m.config.Connections)
}

// findConnectionRecursive recursively searches for a connection
func (m *Manager) findConnectionRecursive(name string, connections []*models.Connection) (*models.Connection, error) {
	for _, conn := range connections {
		if conn.Name == name {
			return conn, nil
		}
		if conn.IsFolder() && len(conn.Children) > 0 {
			if found, err := m.findConnectionRecursive(name, conn.Children); err == nil {
				return found, nil
			}
		}
	}
	return nil, fmt.Errorf("connection '%s' not found", name)
}

// DeleteConnection removes a connection by name
func (m *Manager) DeleteConnection(name string) error {
	if m.config == nil {
		return fmt.Errorf("config not loaded")
	}

	deleted := m.deleteConnectionRecursive(name, &m.config.Connections)
	if !deleted {
		return fmt.Errorf("connection '%s' not found", name)
	}
	return nil
}

// deleteConnectionRecursive recursively deletes a connection
func (m *Manager) deleteConnectionRecursive(name string, connections *[]*models.Connection) bool {
	for i, conn := range *connections {
		if conn.Name == name {
			*connections = append((*connections)[:i], (*connections)[i+1:]...)
			return true
		}
		if conn.IsFolder() && len(conn.Children) > 0 {
			if m.deleteConnectionRecursive(name, &conn.Children) {
				return true
			}
		}
	}
	return false
}

// ListConnections returns all connections in a flat list
func (m *Manager) ListConnections() []*models.Connection {
	if m.config == nil {
		return []*models.Connection{}
	}

	var result []*models.Connection
	m.collectConnectionsRecursive(m.config.Connections, "", &result)
	return result
}

// collectConnectionsRecursive collects all connections with their paths
func (m *Manager) collectConnectionsRecursive(connections []*models.Connection, path string, result *[]*models.Connection) {
	for _, conn := range connections {
		if conn.IsFolder() {
			newPath := path
			if newPath != "" {
				newPath += "/"
			}
			newPath += conn.Name
			m.collectConnectionsRecursive(conn.Children, newPath, result)
		} else {
			*result = append(*result, conn)
		}
	}
}

// findOrCreateFolder finds or creates a folder path
func (m *Manager) findOrCreateFolder(path string) (*models.Connection, error) {
	// Split by both forward slash and backslash
	parts := make([]string, 0)
	current := ""
	for _, ch := range path {
		if ch == '/' || ch == '\\' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}

	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid folder path")
	}

	currentConns := &m.config.Connections
	var parent *models.Connection

	for _, part := range parts {
		found := false
		for _, conn := range *currentConns {
			if conn.Name == part && conn.IsFolder() {
				parent = conn
				currentConns = &conn.Children
				found = true
				break
			}
		}

		if !found {
			// Create the folder
			folder := models.NewFolder(part)
			*currentConns = append(*currentConns, folder)
			parent = folder
			currentConns = &folder.Children
		}
	}

	return parent, nil
}

// UpdateConnection updates an existing connection
func (m *Manager) UpdateConnection(name string, updates *models.Connection) error {
	conn, err := m.FindConnection(name)
	if err != nil {
		return err
	}

	// Update fields if provided
	if updates.Host != "" {
		conn.Host = updates.Host
	}
	if updates.Port != 0 {
		conn.Port = updates.Port
	}
	if updates.Username != "" {
		conn.Username = updates.Username
	}
	if updates.Password != "" {
		conn.Password = updates.Password
	}
	if updates.Domain != "" {
		conn.Domain = updates.Domain
	}
	if updates.Description != "" {
		conn.Description = updates.Description
	}
	if updates.Protocol != "" {
		conn.Protocol = updates.Protocol
	}

	conn.Modified = time.Now().Format(time.RFC3339)
	return nil
}

// encryptPasswords recursively encrypts all passwords in the config
func (m *Manager) encryptPasswords(config *models.Config) {
	m.encryptPasswordsRecursive(config.Connections)
}

// encryptPasswordsRecursive recursively encrypts passwords in connections
func (m *Manager) encryptPasswordsRecursive(connections []*models.Connection) {
	for _, conn := range connections {
		if conn.Password != "" {
			// Don't encrypt 1Password references - keep them as-is
			if !m.onePasswordProvider.IsReference(conn.Password) && !isEncrypted(conn.Password) {
				// Encrypt using DPAPI (Windows) or base64 fallback (other platforms)
				if encrypted, err := encryptPassword(conn.Password); err == nil {
					conn.Password = encrypted
				}
			}
		}
		if len(conn.Children) > 0 {
			m.encryptPasswordsRecursive(conn.Children)
		}
	}
}

// decryptPasswords recursively decrypts all passwords in the config
func (m *Manager) decryptPasswords(config *models.Config) {
	m.decryptPasswordsRecursive(config.Connections)
}

// decryptPasswordsRecursive recursively decrypts passwords in connections
func (m *Manager) decryptPasswordsRecursive(connections []*models.Connection) {
	for _, conn := range connections {
		if conn.Password != "" {
			// Check if it's a 1Password reference first
			if m.onePasswordProvider.IsReference(conn.Password) {
				// Keep 1Password references as-is, don't resolve them during load
				// They will be resolved when actually connecting
				continue
			} else if isEncrypted(conn.Password) {
				// Decrypt using DPAPI (Windows) or base64 fallback (other platforms)
				if decrypted, err := decryptPassword(conn.Password); err == nil {
					conn.Password = decrypted
				} else {
					// If decryption fails (e.g., different user/machine), clear the password
					// User will need to re-enter it
					conn.Password = ""
				}
			}
			// Otherwise it's plain text (legacy or fallback)
		}
		if len(conn.Children) > 0 {
			m.decryptPasswordsRecursive(conn.Children)
		}
	}
}

// isEncrypted checks if a password string is encrypted (base64 encoded)
func isEncrypted(password string) bool {
	// Check if it looks like base64 (length > 20 and only valid base64 chars)
	if len(password) < 20 {
		return false
	}
	for _, ch := range password {
		if !((ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') ||
			(ch >= '0' && ch <= '9') || ch == '+' || ch == '/' || ch == '=') {
			return false
		}
	}
	return true
}

// encryptPassword encrypts a password using Windows DPAPI or fallback
func encryptPassword(password string) (string, error) {
	if runtime.GOOS == "windows" {
		// Use Windows DPAPI via the launcher package
		encrypted, err := launcher.EncryptPasswordWindows(password)
		if err != nil {
			return "", err
		}
		return encrypted, nil
	}
	// For non-Windows, we'd need a different encryption method
	// For now, just return the password (TODO: implement cross-platform encryption)
	return password, nil
}

// decryptPassword decrypts a password using Windows DPAPI or fallback
func decryptPassword(encrypted string) (string, error) {
	if runtime.GOOS == "windows" {
		// Use Windows DPAPI via the launcher package
		decrypted, err := launcher.DecryptPasswordWindows(encrypted)
		if err != nil {
			return "", err
		}
		return decrypted, nil
	}
	// For non-Windows, return as-is
	return encrypted, nil
}

// IsOnePasswordReference checks if a password is a 1Password reference
func (m *Manager) IsOnePasswordReference(password string) bool {
	return m.onePasswordProvider.IsReference(password)
}

// CreateOnePasswordItem creates a new 1Password item and returns the reference
func (m *Manager) CreateOnePasswordItem(vault, title, username, password string) (string, error) {
	return m.onePasswordProvider.CreateItem(vault, title, username, password)
}
