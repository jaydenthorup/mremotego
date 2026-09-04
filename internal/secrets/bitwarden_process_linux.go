//go:build linux

package secrets

import (
	"os/exec"
	"syscall"
)

// setParentDeathSignal makes the kernel send SIGTERM to the child when this
// process exits for any reason.
func setParentDeathSignal(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Pdeathsig = syscall.SIGTERM
}
