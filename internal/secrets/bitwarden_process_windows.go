//go:build windows

package secrets

import (
	"os/exec"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	jobOnce   sync.Once
	jobHandle windows.Handle
)

// killOnCloseJob returns a job object configured to kill its members once the
// last handle to it is closed, which happens when this process dies. That way a
// crash cannot leave an orphaned "bw serve" listening with an unlocked vault.
func killOnCloseJob() windows.Handle {
	jobOnce.Do(func() {
		handle, err := windows.CreateJobObject(nil, nil)
		if err != nil {
			return
		}

		info := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{
			BasicLimitInformation: windows.JOBOBJECT_BASIC_LIMIT_INFORMATION{
				LimitFlags: windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE,
			},
		}
		if _, err := windows.SetInformationJobObject(
			handle,
			windows.JobObjectExtendedLimitInformation,
			uintptr(unsafe.Pointer(&info)),
			uint32(unsafe.Sizeof(info)),
		); err != nil {
			_ = windows.CloseHandle(handle)
			return
		}

		jobHandle = handle
	})

	return jobHandle
}

// configureChildProcess is a no-op on Windows; the console window is already
// hidden by hideConsoleWindow.
func configureChildProcess(cmd *exec.Cmd) {}

// adoptChildProcess assigns the started child to the kill-on-close job object.
// Failures are ignored: losing this safety net is not a reason to refuse the
// connection, and Stop still terminates the process on a clean exit.
func adoptChildProcess(cmd *exec.Cmd) {
	job := killOnCloseJob()
	if job == 0 || cmd.Process == nil {
		return
	}

	handle, err := windows.OpenProcess(
		windows.PROCESS_SET_QUOTA|windows.PROCESS_TERMINATE,
		false,
		uint32(cmd.Process.Pid),
	)
	if err != nil {
		return
	}
	defer windows.CloseHandle(handle)

	_ = windows.AssignProcessToJobObject(job, handle)
}

// stopProcess terminates the child. Windows offers no graceful signal for a
// console-less child process, so it is killed outright.
func stopProcess(cmd *exec.Cmd) {
	if cmd.Process != nil {
		_ = cmd.Process.Kill()
	}
}
