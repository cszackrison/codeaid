package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// RunFirstTimeSetup checks if this is a first-time run and prompts for configuration
// If forceSetup is true, it will run the setup even if a config file exists
func RunFirstTimeSetup(forceSetup bool) error {
	// Check if config file exists
	configPath, err := GetConfigFilePath()
	if err != nil {
		return err
	}

	// If config exists and we're not forcing setup, skip setup
	if !forceSetup {
		if _, err := os.Stat(configPath); err == nil {
			return nil
		}
	}

	// Get existing config to show current values
	existingConfig, _ := Load()
	
	// First time setup
	fmt.Println("\nCodeAid Configuration")
	fmt.Println("==================================================")
	fmt.Println("Press Enter to use default values.")

	// Create a new config with defaults or existing values
	cfg := &Data{
		Model: DefaultModel(),
	}
	
	// If we have an existing config, use its values as defaults
	if existingConfig != nil {
		cfg = existingConfig
	}

	// Get OpenRouter API key and mask it for display
	currentKey := ""
	if cfg.OpenRouterAPIKey != "" {
		// Import utils package here would create circular dependency
		// So implement masking inline
		if len(cfg.OpenRouterAPIKey) > 8 {
			currentKey = cfg.OpenRouterAPIKey[:4] + "..." + cfg.OpenRouterAPIKey[len(cfg.OpenRouterAPIKey)-4:]
		} else {
			currentKey = "****"
		}
		fmt.Printf("Current OpenRouter API Key: %s\n", currentKey)
	}
	
	fmt.Print("OpenRouter API Key: ")
	apiKey := readInput()
	if apiKey != "" {
		cfg.OpenRouterAPIKey = apiKey
	}

	// Model selection
	fmt.Println("\nModel selection:")
	models := AvailableModels

	fmt.Printf("Default model: %s\n\n", cfg.Model)
	for i, model := range models {
		fmt.Printf("%d) %s\n", i+1, model)
	}
	fmt.Printf("%d) Custom model\n", len(models)+1)

	fmt.Print("\nSelect model (1-7): ")
	modelChoice := readInput()
	if modelChoice != "" {
		// Handle custom model
		if modelChoice == fmt.Sprintf("%d", len(models)+1) {
			fmt.Print("Enter custom model identifier: ")
			customModel := readInput()
			if customModel != "" {
				cfg.Model = customModel
			}
		} else {
			// Try to parse as a number and use predefined model
			var idx int
			fmt.Sscanf(modelChoice, "%d", &idx)
			if idx >= 1 && idx <= len(models) {
				cfg.Model = models[idx-1]
			}
		}
	}

	// Save the configuration
	if err := Save(cfg); err != nil {
		return err
	}

	fmt.Printf("\nConfiguration saved to %s\n\n", configPath)
	fmt.Println("Press Enter to continue...")
	readInput()

	return nil
}

// readInput reads a line of input from the user
func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}