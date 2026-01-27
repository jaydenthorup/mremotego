package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

const (
	// EncryptedPrefix is the prefix for encrypted values
	EncryptedPrefix = "enc:"
	// Iterations for PBKDF2
	pbkdf2Iterations = 100000
	// Salt size
	saltSize = 16
	// Key size for AES-256
	keySize = 32
)

// EncryptionProvider handles encryption/decryption of passwords
type EncryptionProvider struct {
	masterPassword string
	enabled        bool
}

// NewEncryptionProvider creates a new encryption provider
func NewEncryptionProvider(masterPassword string) *EncryptionProvider {
	return &EncryptionProvider{
		masterPassword: masterPassword,
		enabled:        masterPassword != "",
	}
}

// IsEnabled returns whether encryption is enabled
func (p *EncryptionProvider) IsEnabled() bool {
	return p.enabled
}

// IsEncrypted checks if a value is encrypted (starts with enc:)
func (p *EncryptionProvider) IsEncrypted(value string) bool {
	return strings.HasPrefix(value, EncryptedPrefix)
}

// Encrypt encrypts a plaintext value using AES-256-GCM
// Returns: "enc:base64(salt+nonce+ciphertext)"
func (p *EncryptionProvider) Encrypt(plaintext string) (string, error) {
	if !p.enabled {
		return "", fmt.Errorf("encryption is not enabled")
	}

	if plaintext == "" {
		return "", nil
	}

	// Generate random salt
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive key from master password using PBKDF2
	key := pbkdf2.Key([]byte(p.masterPassword), salt, pbkdf2Iterations, keySize, sha256.New)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the plaintext
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	// Combine: salt + nonce + ciphertext
	combined := append(salt, nonce...)
	combined = append(combined, ciphertext...)

	// Encode to base64 and add prefix
	encoded := base64.StdEncoding.EncodeToString(combined)
	return EncryptedPrefix + encoded, nil
}

// Decrypt decrypts an encrypted value
// Expects: "enc:base64(salt+nonce+ciphertext)"
func (p *EncryptionProvider) Decrypt(encrypted string) (string, error) {
	if !p.enabled {
		return "", fmt.Errorf("encryption is not enabled")
	}

	if encrypted == "" {
		return "", nil
	}

	if !p.IsEncrypted(encrypted) {
		return "", fmt.Errorf("value is not encrypted")
	}

	// Remove prefix
	encoded := strings.TrimPrefix(encrypted, EncryptedPrefix)

	// Decode from base64
	combined, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("failed to decode encrypted value: %w", err)
	}

	// Extract salt (first 16 bytes)
	if len(combined) < saltSize {
		return "", fmt.Errorf("encrypted value too short")
	}
	salt := combined[:saltSize]
	remaining := combined[saltSize:]

	// Derive key from master password using PBKDF2
	key := pbkdf2.Key([]byte(p.masterPassword), salt, pbkdf2Iterations, keySize, sha256.New)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(remaining) < nonceSize {
		return "", fmt.Errorf("encrypted value too short for nonce")
	}

	// Extract nonce and ciphertext
	nonce := remaining[:nonceSize]
	ciphertext := remaining[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt (wrong password?): %w", err)
	}

	return string(plaintext), nil
}

// DecryptIfNeeded decrypts a value if it's encrypted, otherwise returns it as-is
func (p *EncryptionProvider) DecryptIfNeeded(value string) (string, error) {
	if !p.IsEncrypted(value) {
		return value, nil
	}

	return p.Decrypt(value)
}

// ShouldEncrypt checks if a value should be encrypted
// Returns false for empty strings, 1Password references, or already encrypted values
func (p *EncryptionProvider) ShouldEncrypt(value string) bool {
	if !p.enabled || value == "" {
		return false
	}

	// Don't encrypt if already encrypted
	if p.IsEncrypted(value) {
		return false
	}

	// Don't encrypt 1Password references
	if strings.HasPrefix(value, "op://") {
		return false
	}

	return true
}
