//go:build !windows && !linux

package secrets

import "os/exec"

// setParentDeathSignal is not available on this platform. The server is still
// stopped on a clean exit; only a hard crash can leave it running.
func setParentDeathSignal(cmd *exec.Cmd) {}
