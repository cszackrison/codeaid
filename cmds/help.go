package cmds

import (
	"codeaid/messages"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
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
		var sb strings.Builder
		sb.WriteString("Available commands:\n")

		for _, cmd := range GetAllCommands() {
			sb.WriteString(cmd.Name())
			sb.WriteString(" - ")
			sb.WriteString(cmd.Description())
			sb.WriteString("\n")
		}

		return messages.ResponseMsg(sb.String())
	}
}

// Register the command
func init() {
	RegisterCommand(HelpCommand{})
}
