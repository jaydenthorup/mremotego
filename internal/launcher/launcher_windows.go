//go:build windows
// +build windows

package launcher

import (
	"os/exec"
	"syscall"
)

// hideConsoleWindow sets the command attributes to hide console windows on Windows
func hideConsoleWindow(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
	}
}
