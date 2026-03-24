//go:build windows

package main

import (
	"os/exec"
	"syscall"
)

// configureSysProcAttr configures platform-specific process attributes (Windows)
func configureSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    false,
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
	}
}
