package cmds

import (
	"codeaid/config"
	"codeaid/messages"
	"codeaid/utils"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// ConfigCommand handles configuration updates
type ConfigCommand struct{}

// Name returns the command name
func (c ConfigCommand) Name() string {
	return "/config"
}

// Description returns the command description
func (c ConfigCommand) Description() string {
	return "Update configuration settings"
}

// Execute executes the command
func (c ConfigCommand) Execute(args string) tea.Cmd {
	return func() tea.Msg {
		// Load current configuration
		cfg, err := config.Load()
		if err != nil {
			return messages.CommandResponseMsg(fmt.Sprintf("Error loading configuration: %v", err))
		}

		// Mask API key for display
		maskedKey := utils.MaskAPIKey(cfg.OpenRouterAPIKey)

		// Start config flow with current values
		return messages.ConfigMsg{
			Type:         "init",
			CurrentKey:   maskedKey,
			CurrentModel: cfg.Model,
			PromptText:   fmt.Sprintf("CodeAid Configuration\n====================\nPress Enter to keep current values.\n\nCurrent OpenRouter API Key: %s\nOpenRouter API Key: ", maskedKey),
			ConfigStep:   "api_key",
		}
	}
}