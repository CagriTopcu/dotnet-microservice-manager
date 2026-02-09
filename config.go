package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Service holds information for each microservice
type Service struct {
	Category string `json:"Category"`       // Category: e.g., Category1, Category2
	Name     string `json:"Name"`           // Service name
	Path     string `json:"Path"`           // Service file path
	Port     string `json:"Port,omitempty"` // Port information (optional, can be auto-detected)
}

// Config holds application configuration
type Config struct {
	Services []Service `json:"services"` // List of all services
}

// getConfigPath returns the configuration file path
// Uses %APPDATA%\dotnet-service-manager\config.json on Windows
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

// LoadConfig loads configuration from disk
func LoadConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	// Return empty config if file doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{Services: []Service{}}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse configuration: %v", err)
	}

	return &config, nil
}

// SaveConfig saves configuration to disk
func (c *Config) SaveConfig() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// Save in JSON format (indented)
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to convert config to JSON: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to save configuration: %v", err)
	}

	return nil
}

// ImportServicesFromJSON imports services from a JSON file
func (c *Config) ImportServicesFromJSON(jsonPath string) error {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read JSON file: %v", err)
	}

	var services []Service
	if err := json.Unmarshal(data, &services); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Add new services to existing list (with duplication check)
	for _, newService := range services {
		exists := false
		for _, existingService := range c.Services {
			if existingService.Name == newService.Name && existingService.Category == newService.Category {
				exists = true
				break
			}
		}
		if !exists {
			c.Services = append(c.Services, newService)
		}
	}

	return nil
}

// GetServicesByCategory returns services in a specific category
func (c *Config) GetServicesByCategory(category string) []Service {
	var result []Service
	for _, service := range c.Services {
		if service.Category == category {
			result = append(result, service)
		}
	}
	return result
}

// GetCategories returns all categories (unique)
func (c *Config) GetCategories() []string {
	categoryMap := make(map[string]bool)
	for _, service := range c.Services {
		categoryMap[service.Category] = true
	}

	var categories []string
	for category := range categoryMap {
		categories = append(categories, category)
	}
	return categories
}

// ClearAllServices clears all services
func (c *Config) ClearAllServices() {
	c.Services = []Service{}
}

// GetServiceCount returns total service count
func (c *Config) GetServiceCount() int {
	return len(c.Services)
}
