//go:build windows

package launcher

import (
	"encoding/binary"
	"errors"
	"unicode/utf16"

	"github.com/danieljoos/wincred"
	"golang.org/x/sys/windows"
)

// writeGenericCredential stores a credential in the Windows Credential Manager.
//
// This replaces shelling out to cmdkey, which takes the password on its command
// line where every local process can read it.
//
// The blob must be UTF-16LE: that is what cmdkey writes and what mstsc expects.
// Storing plain UTF-8 makes the credential look valid while silently failing to
// log in. Persistence is enterprise, matching "cmdkey /generic".
func writeGenericCredential(target, username, password string) error {
	cred := wincred.NewGenericCredential(target)
	cred.UserName = username
	cred.CredentialBlob = utf16LE(password)
	cred.Persist = wincred.PersistEnterprise

	return cred.Write()
}

// deleteGenericCredential removes a credential. A credential that is already
// gone is not an error.
func deleteGenericCredential(target string) error {
	cred, err := wincred.GetGenericCredential(target)
	if err != nil {
		if errors.Is(err, windows.ERROR_NOT_FOUND) {
			return nil
		}
		return err
	}

	return cred.Delete()
}

// utf16LE encodes a string as UTF-16 little endian without a byte order mark
// and without a terminating NUL, which is how Windows stores credential blobs.
func utf16LE(s string) []byte {
	encoded := utf16.Encode([]rune(s))

	out := make([]byte, len(encoded)*2)
	for i, r := range encoded {
		binary.LittleEndian.PutUint16(out[i*2:], r)
	}

	return out
}
