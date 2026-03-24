//go:build !windows

package main

import (
	"os"
	"os/exec"
	"syscall"
	"time"
)

// configureSysProcAttr configures platform-specific process attributes (Unix/Linux)
// On Unix systems, we create a process group so we can kill all child processes
func configureSysProcAttr(cmd *exec.Cmd) {
	// Create a new process group for the command
	// This allows us to terminate the entire process tree
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0, // Create new process group
	}
}

// terminateProcess stops a process gracefully on Unix/Linux
// First tries SIGTERM (graceful), then SIGKILL if needed
func terminateProcess(proc *os.Process) error {
	// First, try graceful shutdown with SIGTERM
	if err := syscall.Kill(-proc.Pid, syscall.SIGTERM); err != nil {
		// If SIGTERM fails, fall back to direct Kill
		return proc.Kill()
	}

	// Give the process time to clean up
	time.Sleep(500 * time.Millisecond)

	// Check if process is still running by trying to signal it (0 means check if running)
	if err := syscall.Kill(-proc.Pid, 0); err == nil {
		// Process is still running, force kill the entire process group
		if err := syscall.Kill(-proc.Pid, syscall.SIGKILL); err != nil {
			// Process already terminated
			return nil
		}
	}

	return nil
}
