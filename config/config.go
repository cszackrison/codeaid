package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Data represents the application configuration
type Data struct {
	OpenRouterAPIKey string `json:"openrouter_api_key"`
	Model            string `json:"model"`
}

// DefaultModel returns the default model identifier
func DefaultModel() string {
	return "mistralai/mistral-small-3.1-24b-instruct:free"
}

// GetConfigDir returns the configuration directory path
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".config", "codeaid")
	return configDir, nil
}

// GetConfigFilePath returns the path to the config file
func GetConfigFilePath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "config.json"), nil
}

// Load loads the configuration from disk
func Load() (*Data, error) {
	configPath, err := GetConfigFilePath()
	if err != nil {
		return nil, err
	}

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return default config if file doesn't exist
		return &Data{
			Model: DefaultModel(),
		}, nil
	}

	// Read and parse config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Data
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Save saves the configuration to disk
func Save(config *Data) error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configPath, err := GetConfigFilePath()
	if err != nil {
		return err
	}

	// Write config file
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return err
	}

	return nil
}