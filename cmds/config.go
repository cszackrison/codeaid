package cmds

import (
	"codeaid/config"
	"codeaid/messages"
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
		maskedKey := ""
		if cfg.OpenRouterAPIKey != "" {
			if len(cfg.OpenRouterAPIKey) > 8 {
				maskedKey = cfg.OpenRouterAPIKey[:4] + "..." + cfg.OpenRouterAPIKey[len(cfg.OpenRouterAPIKey)-4:]
			} else {
				maskedKey = "****"
			}
		}

		// Start config flow with current values
		return messages.ConfigInitMsg{
			CurrentAPIKey: maskedKey,
			CurrentModel:  cfg.Model,
		}
	}
}