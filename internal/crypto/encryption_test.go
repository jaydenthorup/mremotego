package crypto

import "testing"

func TestShouldEncrypt(t *testing.T) {
	provider := NewEncryptionProvider("master")

	tests := map[string]bool{
		"hunter2":                      true,
		"":                             false,
		"op://Private/Server/password": false,
		"bw://8f3c-item-id":            false,
		"enc:AAAA":                     false,
	}

	for value, want := range tests {
		if got := provider.ShouldEncrypt(value); got != want {
			t.Errorf("ShouldEncrypt(%q) = %v, want %v", value, got, want)
		}
	}
}

func TestShouldEncryptDisabledWithoutMasterPassword(t *testing.T) {
	provider := NewEncryptionProvider("")

	if provider.ShouldEncrypt("hunter2") {
		t.Error("ShouldEncrypt should be false when no master password is set")
	}
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	provider := NewEncryptionProvider("master")

	encrypted, err := provider.Encrypt("hunter2")
	if err != nil {
		t.Fatalf("Encrypt returned error: %v", err)
	}
	if !provider.IsEncrypted(encrypted) {
		t.Fatalf("Encrypt produced %q which is not recognised as encrypted", encrypted)
	}

	decrypted, err := provider.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt returned error: %v", err)
	}
	if decrypted != "hunter2" {
		t.Errorf("Decrypt = %q, want %q", decrypted, "hunter2")
	}
}

func TestDecryptWithWrongPasswordFails(t *testing.T) {
	encrypted, err := NewEncryptionProvider("master").Encrypt("hunter2")
	if err != nil {
		t.Fatalf("Encrypt returned error: %v", err)
	}

	if _, err := NewEncryptionProvider("other").Decrypt(encrypted); err == nil {
		t.Error("expected decryption with the wrong master password to fail")
	}
}
