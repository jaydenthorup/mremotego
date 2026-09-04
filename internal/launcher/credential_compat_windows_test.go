//go:build windows

package launcher

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/danieljoos/wincred"
)

// TestCredentialMatchesCmdkey verifies that the credential written through the
// Credential Manager API is byte for byte what "cmdkey /generic" would have
// written. mstsc reads the blob as UTF-16LE, so storing plain UTF-8 would
// produce a credential that looks valid but silently fails to log in.
func TestCredentialMatchesCmdkey(t *testing.T) {
	if _, err := exec.LookPath("cmdkey"); err != nil {
		t.Skip("cmdkey not available")
	}

	const (
		user     = "tester"
		password = "passw0rd123"
		refTgt   = "TERMSRV/mremotego-selftest-cmdkey"
		apiTgt   = "TERMSRV/mremotego-selftest-api"
	)

	// Start from a clean slate: a leftover credential from an earlier run would
	// make the comparison meaningless.
	_ = deleteGenericCredential(refTgt)
	_ = deleteGenericCredential(apiTgt)

	out, err := exec.Command("cmdkey", "/generic:"+refTgt, "/user:"+user, "/pass:"+password).CombinedOutput()
	if err != nil {
		t.Skipf("cmdkey refused to store the reference credential: %v (%s)", err, out)
	}
	defer deleteGenericCredential(refTgt)

	reference, err := wincred.GetGenericCredential(refTgt)
	if err != nil {
		t.Skipf("cmdkey credential could not be read back: %v", err)
	}

	if err := writeGenericCredential(apiTgt, user, password); err != nil {
		t.Fatalf("writeGenericCredential returned error: %v", err)
	}
	defer deleteGenericCredential(apiTgt)

	written, err := wincred.GetGenericCredential(apiTgt)
	if err != nil {
		t.Fatalf("failed to read back the credential we wrote: %v", err)
	}

	if !bytes.Equal(reference.CredentialBlob, written.CredentialBlob) {
		t.Errorf("blob mismatch:\n cmdkey % x\n api    % x",
			reference.CredentialBlob, written.CredentialBlob)
	}
	if written.UserName != reference.UserName {
		t.Errorf("user name = %q, want %q", written.UserName, reference.UserName)
	}
	if written.Persist != reference.Persist {
		t.Errorf("persist = %v, want %v", written.Persist, reference.Persist)
	}
}

// TestUTF16LE pins the encoding independently of whether cmdkey is available.
func TestUTF16LE(t *testing.T) {
	got := utf16LE("ab")
	want := []byte{'a', 0x00, 'b', 0x00}

	if !bytes.Equal(got, want) {
		t.Errorf("utf16LE(\"ab\") = % x, want % x", got, want)
	}

	// No byte order mark and no terminating NUL, matching what cmdkey stores.
	if len(utf16LE("")) != 0 {
		t.Errorf("utf16LE(\"\") = % x, want empty", utf16LE(""))
	}
}

// TestDeleteMissingCredentialIsNotAnError guards the cleanup path, which runs
// for every RDP connection that is removed.
func TestDeleteMissingCredentialIsNotAnError(t *testing.T) {
	if err := deleteGenericCredential("TERMSRV/mremotego-selftest-absent"); err != nil {
		t.Errorf("deleting a missing credential returned: %v", err)
	}
}
