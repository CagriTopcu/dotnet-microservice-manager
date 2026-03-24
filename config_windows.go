//go:build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// getConfigPath returns the configuration file path (Windows)
// Uses %APPDATA%\dotnet-service-manager\config.json
func getConfigPath() (string, error) {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		return "", fmt.Errorf("APPDATA environment variable not found")
	}

	configDir := filepath.Join(appData, "dotnet-service-manager")

	// Create folder if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %v", err)
	}

	return filepath.Join(configDir, "config.json"), nil
}
