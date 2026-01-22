//go:build !windows
// +build !windows

package launcher

// encryptPasswordWindows is a no-op on non-Windows platforms
func encryptPasswordWindows(password string) (string, error) {
	return "", nil
}

// decryptPasswordWindows is a no-op on non-Windows platforms
func decryptPasswordWindows(encryptedPassword string) (string, error) {
	return "", nil
}

// EncryptPasswordWindows is an exported wrapper (no-op on non-Windows)
func EncryptPasswordWindows(password string) (string, error) {
	return password, nil // Just return plaintext on non-Windows
}

// DecryptPasswordWindows is an exported wrapper (no-op on non-Windows)
func DecryptPasswordWindows(encrypted string) (string, error) {
	return encrypted, nil // Just return as-is on non-Windows
}
