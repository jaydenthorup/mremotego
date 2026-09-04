package launcher

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	// passwordFilePrefix marks the temporary files holding a password.
	passwordFilePrefix = "pw-"
	// passwordFileMaxAge is how long a leftover file may survive before it is
	// considered stale and removed at start-up.
	passwordFileMaxAge = 2 * time.Minute
	// passwordFileGrace bounds how long the file may exist after the helper
	// program was started, in case that program never exits.
	passwordFileGrace = 5 * time.Second
	// staleCleanupDelay is the safety net for terminal emulators that never ran
	// the generated snippet, which would otherwise leave the file behind.
	staleCleanupDelay = 60 * time.Second
)

// tempDir is the directory used for the short lived files handed to helper
// programs. It is the same directory the generated .rdp files already use.
func tempDir() string {
	return filepath.Join(os.TempDir(), "mremotego")
}

// writePasswordFile stores a password in a private temporary file and returns
// its path together with a cleanup function.
//
// Passing a password through a file rather than the command line keeps it out
// of the process list, where any local user could read it.
//
// os.CreateTemp creates the file with mode 0600, and on Windows the per-user
// temporary directory is already restricted to the user, SYSTEM and
// administrators - the same protection the generated .rdp files rely on. The
// window in which the file exists is short: the helper program reads it while
// starting up and the caller removes it right after.
func writePasswordFile(password string) (path string, cleanup func(), err error) {
	dir := tempDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", func() {}, fmt.Errorf("failed to create temp directory: %w", err)
	}

	file, err := os.CreateTemp(dir, passwordFilePrefix+"*.txt")
	if err != nil {
		return "", func() {}, fmt.Errorf("failed to create password file: %w", err)
	}
	path = file.Name()

	// PuTTY and the sshpass wrapper both read the first line of the file.
	if _, err := file.WriteString(password + "\n"); err != nil {
		file.Close()
		os.Remove(path)
		return "", func() {}, fmt.Errorf("failed to write password file: %w", err)
	}

	if err := file.Close(); err != nil {
		os.Remove(path)
		return "", func() {}, fmt.Errorf("failed to close password file: %w", err)
	}

	var once sync.Once
	cleanup = func() {
		once.Do(func() {
			// Overwrite before unlinking so the content does not linger in
			// slack space if the delete fails.
			if f, err := os.OpenFile(path, os.O_WRONLY, 0600); err == nil {
				_, _ = f.Write(make([]byte, len(password)+1))
				f.Close()
			}
			os.Remove(path)
		})
	}

	return path, cleanup, nil
}

// cleanupStalePasswordFiles removes password files left behind by an earlier
// run that was killed before it could clean up.
func cleanupStalePasswordFiles() {
	entries, err := os.ReadDir(tempDir())
	if err != nil {
		return
	}

	cutoff := time.Now().Add(-passwordFileMaxAge)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasPrefix(entry.Name(), passwordFilePrefix) {
			continue
		}

		info, err := entry.Info()
		if err != nil || info.ModTime().After(cutoff) {
			continue
		}

		os.Remove(filepath.Join(tempDir(), entry.Name()))
	}
}
