//go:build !windows

package main

import (
	"os/exec"
)

// configureSysProcAttr configures platform-specific process attributes (Unix/Linux)
// On Unix systems, we don't need special process attributes for this use case
func configureSysProcAttr(cmd *exec.Cmd) {
	// No platform-specific configuration needed for Unix/Linux
	cmd.SysProcAttr = nil
}
