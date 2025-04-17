package cmds

import (
	"codeaid/messages"
	tea "github.com/charmbracelet/bubbletea"
)

// HelpCommand displays help information
type HelpCommand struct{}

// Name returns the command name
func (c HelpCommand) Name() string {
	return "/help"
}

// Description returns the command description
func (c HelpCommand) Description() string {
	return "Show available commands"
}

// Execute executes the command
func (c HelpCommand) Execute(args string) tea.Cmd {
	return func() tea.Msg {
		// Collect all commands
		allCommands := GetAllCommands()
		cmdInfos := make([]messages.CommandInfo, 0, len(allCommands))
		
		for _, cmd := range allCommands {
			cmdInfos = append(cmdInfos, messages.CommandInfo{
				Name:        cmd.Name(),
				Description: cmd.Description(),
			})
		}
		
		// Return a specialized help message that will be styled in main.go
		return messages.HelpMsg{
			Header:   "Available commands:",
			Commands: cmdInfos,
		}
	}
}

// Register the command
func init() {
	RegisterCommand(HelpCommand{})
}