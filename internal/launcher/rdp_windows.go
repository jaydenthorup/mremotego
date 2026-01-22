//go:build windows
// +build windows

package launcher

import (
	"encoding/base64"
	"fmt"
	"syscall"
	"unsafe"
)

var (
	crypt32                = syscall.NewLazyDLL("crypt32.dll")
	procCryptProtectData   = crypt32.NewProc("CryptProtectData")
	procCryptUnprotectData = crypt32.NewProc("CryptUnprotectData")
)

type dataBlob struct {
	cbData uint32
	pbData *byte
}

// encryptPasswordWindows encrypts a password using Windows DPAPI (CryptProtectData)
// This is the same method mRemoteNG uses for RDP password encryption
func encryptPasswordWindows(password string) (string, error) {
	if password == "" {
		return "", nil
	}

	// Convert password to UTF-16LE (Windows unicode) as mstsc expects
	passwordBytes := encodeUTF16LE(password)

	var inBlob dataBlob
	inBlob.cbData = uint32(len(passwordBytes))
	inBlob.pbData = &passwordBytes[0]

	var outBlob dataBlob

	// Call CryptProtectData with CRYPTPROTECT_UI_FORBIDDEN flag
	// BOOL CryptProtectData(
	//   DATA_BLOB *pDataIn,
	//   LPCWSTR szDataDescr,
	//   DATA_BLOB *pOptionalEntropy,
	//   PVOID pvReserved,
	//   CRYPTPROTECT_PROMPTSTRUCT *pPromptStruct,
	//   DWORD dwFlags,
	//   DATA_BLOB *pDataOut
	// );
	const CRYPTPROTECT_UI_FORBIDDEN = 0x1
	ret, _, err := procCryptProtectData.Call(
		uintptr(unsafe.Pointer(&inBlob)),
		0,                         // szDataDescr
		0,                         // pOptionalEntropy
		0,                         // pvReserved
		0,                         // pPromptStruct
		CRYPTPROTECT_UI_FORBIDDEN, // dwFlags - suppress UI
		uintptr(unsafe.Pointer(&outBlob)),
	)

	if ret == 0 {
		return "", fmt.Errorf("CryptProtectData failed: %v", err)
	}

	// Convert encrypted data to byte slice
	encryptedData := make([]byte, outBlob.cbData)
	for i := uint32(0); i < outBlob.cbData; i++ {
		encryptedData[i] = *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(outBlob.pbData)) + uintptr(i)))
	}

	// Free the output buffer
	syscall.LocalFree(syscall.Handle(unsafe.Pointer(outBlob.pbData)))

	// Encode as base64
	encoded := base64.StdEncoding.EncodeToString(encryptedData)

	return encoded, nil
}

// decryptPasswordWindows decrypts a password using Windows DPAPI (CryptUnprotectData)
func decryptPasswordWindows(encryptedPassword string) (string, error) {
	if encryptedPassword == "" {
		return "", nil
	}

	// Decode base64
	encryptedData, err := base64.StdEncoding.DecodeString(encryptedPassword)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	var inBlob dataBlob
	inBlob.cbData = uint32(len(encryptedData))
	inBlob.pbData = &encryptedData[0]

	var outBlob dataBlob

	// Call CryptUnprotectData
	ret, _, err := procCryptUnprotectData.Call(
		uintptr(unsafe.Pointer(&inBlob)),
		0, // ppszDataDescr
		0, // pOptionalEntropy
		0, // pvReserved
		0, // pPromptStruct
		0, // dwFlags
		uintptr(unsafe.Pointer(&outBlob)),
	)

	if ret == 0 {
		return "", fmt.Errorf("CryptUnprotectData failed: %v", err)
	}

	// Convert decrypted data to string
	decryptedData := make([]byte, outBlob.cbData)
	for i := uint32(0); i < outBlob.cbData; i++ {
		decryptedData[i] = *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(outBlob.pbData)) + uintptr(i)))
	}

	// Free the output buffer
	syscall.LocalFree(syscall.Handle(unsafe.Pointer(outBlob.pbData)))

	return string(decryptedData), nil
}

// encodeUTF16LE encodes a string to UTF-16 Little Endian (required for Windows RDP passwords)
func encodeUTF16LE(s string) []byte {
	runes := []rune(s)
	result := make([]byte, 0, len(runes)*2+2) // +2 for null terminator

	for _, r := range runes {
		// Convert rune to UTF-16LE bytes
		result = append(result, byte(r), byte(r>>8))
	}

	// Add null terminator
	result = append(result, 0, 0)

	return result
}

// EncryptPasswordWindows is an exported wrapper for password encryption
func EncryptPasswordWindows(password string) (string, error) {
	return encryptPasswordWindows(password)
}

// DecryptPasswordWindows is an exported wrapper for password decryption
func DecryptPasswordWindows(encrypted string) (string, error) {
	return decryptPasswordWindows(encrypted)
}
