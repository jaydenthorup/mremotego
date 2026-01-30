package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jaydenthorup/mremotego/internal/crypto"
	"github.com/jaydenthorup/mremotego/internal/secrets"
	"github.com/jaydenthorup/mremotego/pkg/models"
	"gopkg.in/yaml.v3"
)

// Manager handles configuration file operations
type Manager struct {
	configPath             string
	config                 *models.Config
	onePasswordProvider    *secrets.OnePasswordProvider    // CLI provider (fallback)
	onePasswordSDKProvider *secrets.OnePasswordSDKProvider // SDK provider (preferred)
	encryptionProvider     *crypto.EncryptionProvider
}

// NewManager creates a new configuration manager
func NewManager(configPath string) *Manager {
	return &Manager{
		configPath:          configPath,
		onePasswordProvider: secrets.NewOnePasswordProvider(),
		encryptionProvider:  nil, // Will be set when master password is provided
	}
}

// SetMasterPassword sets the master password for encryption/decryption
func (m *Manager) SetMasterPassword(password string) {
	m.encryptionProvider = crypto.NewEncryptionProvider(password)
}

// GetConfigPath returns the current config file path
func (m *Manager) GetConfigPath() string {
	return m.configPath
}

// GetDefaultConfigPath returns the default configuration file path
// It checks for a recent file first, then falls back to the default location
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

	// Check if there's a recent file saved
	recentFilePath := filepath.Join(configDir, "recent.txt")
	if data, err := os.ReadFile(recentFilePath); err == nil {
		recentPath := string(data)
		// Verify the file still exists
		if _, err := os.Stat(recentPath); err == nil {
			return recentPath, nil
		}
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

	// Debug: Print what was loaded
	fmt.Printf("[DEBUG] Config loaded - Version: %s\n", config.Version)
	fmt.Printf("[DEBUG] Settings pointer: %v\n", config.Settings)
	if config.Settings != nil {
		fmt.Printf("[DEBUG] OnePasswordAccount: '%s'\n", config.Settings.OnePasswordAccount)
	}

	// Ensure settings struct is always initialized (even for older config files)
	if config.Settings == nil {
		config.Settings = &models.Settings{}
		fmt.Println("[DEBUG] Initialized nil Settings struct")
	}

	// Decrypt passwords if encryption is enabled
	if m.encryptionProvider != nil && m.encryptionProvider.IsEnabled() {
		if err := m.decryptPasswords(&config); err != nil {
			return fmt.Errorf("failed to decrypt passwords: %w", err)
		}
	}

	m.config = &config

	// Save this as the most recently used config file
	m.saveRecentFile()

	return nil
}

// decryptPasswords recursively decrypts all encrypted passwords in the config
func (m *Manager) decryptPasswords(config *models.Config) error {
	return m.decryptPasswordsRecursive(config.Connections)
}

func (m *Manager) decryptPasswordsRecursive(connections []*models.Connection) error {
	for _, conn := range connections {
		if conn.Password != "" && m.encryptionProvider.IsEncrypted(conn.Password) {
			decrypted, err := m.encryptionProvider.Decrypt(conn.Password)
			if err != nil {
				return fmt.Errorf("failed to decrypt password for '%s': %w", conn.Name, err)
			}
			conn.Password = decrypted
		}

		// Recursively decrypt children
		if conn.IsFolder() && len(conn.Children) > 0 {
			if err := m.decryptPasswordsRecursive(conn.Children); err != nil {
				return err
			}
		}
	}
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

	// Create a copy for encryption (don't modify the in-memory config)
	configCopy := m.config.DeepCopy()

	// Encrypt passwords if encryption is enabled
	if m.encryptionProvider != nil && m.encryptionProvider.IsEnabled() {
		if err := m.encryptPasswords(configCopy); err != nil {
			return fmt.Errorf("failed to encrypt passwords: %w", err)
		}
	}

	data, err := yaml.Marshal(configCopy)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(m.configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// encryptPasswords recursively encrypts all passwords that should be encrypted
func (m *Manager) encryptPasswords(config *models.Config) error {
	return m.encryptPasswordsRecursive(config.Connections)
}

func (m *Manager) encryptPasswordsRecursive(connections []*models.Connection) error {
	for _, conn := range connections {
		if m.encryptionProvider.ShouldEncrypt(conn.Password) {
			encrypted, err := m.encryptionProvider.Encrypt(conn.Password)
			if err != nil {
				return fmt.Errorf("failed to encrypt password for '%s': %w", conn.Name, err)
			}
			conn.Password = encrypted
		}

		// Recursively encrypt children
		if conn.IsFolder() && len(conn.Children) > 0 {
			if err := m.encryptPasswordsRecursive(conn.Children); err != nil {
				return err
			}
		}
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

// IsOnePasswordReference checks if a password is a 1Password reference
func (m *Manager) IsOnePasswordReference(password string) bool {
	return m.onePasswordProvider.IsReference(password)
}

// CreateOnePasswordItem creates a new 1Password item and returns the reference
func (m *Manager) CreateOnePasswordItem(vault, title, username, password string) (string, error) {
	// Use SDK provider if available (preferred), otherwise fall back to CLI
	if m.onePasswordSDKProvider != nil {
		return m.onePasswordSDKProvider.CreateItem(vault, title, username, password)
	}
	return m.onePasswordProvider.CreateItem(vault, title, username, password)
}

// SetOnePasswordSDKProvider sets the SDK provider for 1Password operations
// Also loads vault name mappings from settings if available
func (m *Manager) SetOnePasswordSDKProvider(provider *secrets.OnePasswordSDKProvider) {
	m.onePasswordSDKProvider = provider
	
	// Load vault name mappings from config settings if available
	if m.config != nil && m.config.Settings != nil && len(m.config.Settings.VaultNames) > 0 {
		provider.SetVaultNameMap(m.config.Settings.VaultNames)
		fmt.Printf("[Config] Loaded %d vault name mappings from settings\n", len(m.config.Settings.VaultNames))
	}
}

// SaveVaultNameMappings saves vault name mappings to the config settings
func (m *Manager) SaveVaultNameMappings(mappings map[string]string) error {
	if m.config == nil {
		return fmt.Errorf("config not loaded")
	}
	
	if m.config.Settings == nil {
		m.config.Settings = &models.Settings{}
	}
	
	m.config.Settings.VaultNames = mappings
	
	// Also update the SDK provider if it's set
	if m.onePasswordSDKProvider != nil {
		m.onePasswordSDKProvider.SetVaultNameMap(mappings)
	}
	
	// Save the config
	return m.Save()
}

// saveRecentFile saves the current config path as the most recently used file
func (m *Manager) saveRecentFile() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	var configDir string
	if os.Getenv("APPDATA") != "" {
		configDir = filepath.Join(os.Getenv("APPDATA"), "mremotego")
	} else {
		configDir = filepath.Join(homeDir, ".config", "mremotego")
	}

	// Ensure directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	recentFilePath := filepath.Join(configDir, "recent.txt")
	// Get absolute path
	absPath, err := filepath.Abs(m.configPath)
	if err != nil {
		absPath = m.configPath
	}

	return os.WriteFile(recentFilePath, []byte(absPath), 0644)
}
