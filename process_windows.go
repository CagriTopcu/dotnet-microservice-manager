//go:build windows

package main

import (
	"os"
	"os/exec"
	"syscall"
	"time"
)

// configureSysProcAttr configures platform-specific process attributes (Windows)
func configureSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    false,
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
	}
}

// terminateProcess stops a process (Windows)
// On Windows, Kill() will terminate the process and its children
func terminateProcess(proc *os.Process) error {
	return proc.Kill()
}
