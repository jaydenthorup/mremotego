package launcher

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestWritePasswordFile(t *testing.T) {
	path, cleanup, err := writePasswordFile("hunter2")
	if err != nil {
		t.Fatalf("writePasswordFile returned error: %v", err)
	}
	defer cleanup()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read password file: %v", err)
	}

	// PuTTY and the sshpass wrapper read the first line of the file.
	if string(data) != "hunter2\n" {
		t.Errorf("password file contains %q, want %q", string(data), "hunter2\n")
	}
}

func TestWritePasswordFileCleanupRemovesFile(t *testing.T) {
	path, cleanup, err := writePasswordFile("hunter2")
	if err != nil {
		t.Fatalf("writePasswordFile returned error: %v", err)
	}

	cleanup()

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("password file still exists after cleanup (stat error: %v)", err)
	}
}

func TestWritePasswordFileCleanupIsIdempotent(t *testing.T) {
	_, cleanup, err := writePasswordFile("hunter2")
	if err != nil {
		t.Fatalf("writePasswordFile returned error: %v", err)
	}

	cleanup()
	cleanup() // must not panic or fail
}

func TestWritePasswordFileIsUnique(t *testing.T) {
	first, cleanupFirst, err := writePasswordFile("a")
	if err != nil {
		t.Fatalf("writePasswordFile returned error: %v", err)
	}
	defer cleanupFirst()

	second, cleanupSecond, err := writePasswordFile("b")
	if err != nil {
		t.Fatalf("writePasswordFile returned error: %v", err)
	}
	defer cleanupSecond()

	if first == second {
		t.Error("two password files got the same path")
	}
}

func TestCleanupStalePasswordFiles(t *testing.T) {
	stale := filepath.Join(tempDir(), passwordFilePrefix+"stale-test.txt")
	if err := os.MkdirAll(tempDir(), 0700); err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	if err := os.WriteFile(stale, []byte("old\n"), 0600); err != nil {
		t.Fatalf("failed to create stale file: %v", err)
	}

	old := time.Now().Add(-2 * passwordFileMaxAge)
	if err := os.Chtimes(stale, old, old); err != nil {
		t.Fatalf("failed to age the stale file: %v", err)
	}

	fresh, cleanupFresh, err := writePasswordFile("keep me")
	if err != nil {
		t.Fatalf("writePasswordFile returned error: %v", err)
	}
	defer cleanupFresh()

	cleanupStalePasswordFiles()

	if _, err := os.Stat(stale); !os.IsNotExist(err) {
		os.Remove(stale)
		t.Errorf("stale password file survived cleanup (stat error: %v)", err)
	}
	if _, err := os.Stat(fresh); err != nil {
		t.Errorf("recent password file was removed: %v", err)
	}
}

func TestShellCommandQuotesArguments(t *testing.T) {
	got := shellCommand("ssh", "-p", "22", "user@host")
	want := "'ssh' '-p' '22' 'user@host'"

	if got != want {
		t.Errorf("shellCommand = %q, want %q", got, want)
	}
}

func TestShellCommandEscapesSingleQuotes(t *testing.T) {
	got := shellCommand("ssh", "o'brien@host")

	if strings.Contains(got, "o'brien@host") {
		t.Errorf("shellCommand left an unescaped quote: %q", got)
	}
	if !strings.Contains(got, `'\''`) {
		t.Errorf("shellCommand = %q, want the single quote shell-escaped", got)
	}
}

// The password must reach sshpass through the environment, never through argv,
// where every local user could read it.
func TestSshpassCommandKeepsPasswordOutOfArgv(t *testing.T) {
	got := sshpassCommand("/tmp/mremotego/pw-1.txt", []string{"-p", "22", "user@host"})

	if !strings.Contains(got, "SSHPASS=") {
		t.Errorf("snippet does not set SSHPASS: %q", got)
	}
	if !strings.Contains(got, "'sshpass' '-e' 'ssh'") {
		t.Errorf("snippet does not run sshpass -e ssh: %q", got)
	}
	if strings.Contains(got, "-p '") && strings.Contains(got, "sshpass' '-p") {
		t.Errorf("snippet passes the password on the command line: %q", got)
	}
	if !strings.Contains(got, "rm -f") {
		t.Errorf("snippet does not remove the password file: %q", got)
	}
}
