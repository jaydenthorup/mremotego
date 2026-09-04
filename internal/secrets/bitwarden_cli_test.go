package secrets

import (
	"errors"
	"os/exec"
	"testing"
)

func TestClassifyStartupFailure(t *testing.T) {
	tests := []struct {
		output string
		want   error
	}{
		{output: "You are not logged in.", want: ErrNotAuthenticated},
		{output: "Vault is locked.", want: ErrVaultLocked},
		{output: "EADDRINUSE: address already in use", want: nil},
		{output: "", want: nil},
	}

	for _, tt := range tests {
		got := classifyStartupFailure(tt.output)
		if tt.want == nil {
			if got != nil {
				t.Errorf("classifyStartupFailure(%q) = %v, want nil", tt.output, got)
			}
			continue
		}
		if !errors.Is(got, tt.want) {
			t.Errorf("classifyStartupFailure(%q) = %v, want %v", tt.output, got, tt.want)
		}
	}
}

// TestRealBwServeLifecycle exercises the actual Bitwarden CLI when it is
// installed. It does not need a logged-in vault: an unauthenticated CLI is a
// valid outcome and is checked for the right error.
func TestRealBwServeLifecycle(t *testing.T) {
	if _, err := exec.LookPath("bw"); err != nil {
		t.Skip("bw CLI not in PATH")
	}

	provider := NewBitwardenProvider()
	defer provider.Close()

	if !provider.IsEnabled() {
		t.Fatal("provider reports disabled although bw is in PATH")
	}

	status, err := provider.Status()
	switch {
	case errors.Is(err, ErrNotAuthenticated):
		// The CLI refuses to serve without a login, which is the expected
		// result on a machine that has never run "bw login".
		t.Log("bw is not logged in; start-up failure correctly reported")
		return
	case err != nil:
		t.Fatalf("Status returned error: %v", err)
	}

	t.Logf("vault status: %q", status)

	provider.mu.Lock()
	server := provider.server
	provider.mu.Unlock()

	if server == nil {
		t.Fatal("no bw serve process was started")
	}
	t.Logf("bw serve pid=%d port=%d", server.cmd.Process.Pid, server.port)

	if err := provider.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}
	if !server.Exited() {
		t.Error("bw serve is still running after Close")
	}
}
