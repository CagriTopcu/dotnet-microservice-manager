//go:build !windows

package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// getConfigPath returns the configuration file path (Unix/Linux)
// Uses ~/.config/dotnet-service-manager/config.json (XDG Base Directory)
func getConfigPath() (string, error) {
	// Try XDG_CONFIG_HOME first
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		// Fall back to ~/.config
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to determine home directory: %v", err)
		}
		configHome = filepath.Join(home, ".config")
	}

	configDir := filepath.Join(configHome, "dotnet-service-manager")

	// Create folder if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %v", err)
	}

	return filepath.Join(configDir, "config.json"), nil
}
