//go:build !windows

package secrets

import (
	"os/exec"
	"syscall"
)

// configureChildProcess asks the kernel to signal the child when this process
// dies, so a crash cannot leave an orphaned "bw serve" running.
func configureChildProcess(cmd *exec.Cmd) {
	setParentDeathSignal(cmd)
}

// adoptChildProcess has nothing to do on Unix; the death signal is configured
// before the process starts.
func adoptChildProcess(cmd *exec.Cmd) {}

// stopProcess asks the child to terminate.
func stopProcess(cmd *exec.Cmd) {
	if cmd.Process == nil {
		return
	}
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		_ = cmd.Process.Kill()
	}
}
