package cmds

import (
	tea "github.com/charmbracelet/bubbletea"
)

// ExitCommand exits the application
type ExitCommand struct{}

// Name returns the command name
func (c ExitCommand) Name() string {
	return "/exit"
}

// Description returns the command description
func (c ExitCommand) Description() string {
	return "Exit the application"
}

// Execute executes the command
func (c ExitCommand) Execute(args string) tea.Cmd {
	return tea.Quit
}

// Register the command
func init() {
	RegisterCommand(ExitCommand{})
}
