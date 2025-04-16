package cmds

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Command defines the interface for all commands
type Command interface {
	// Name returns the command name including the leading slash
	Name() string

	// Description returns a brief description of what the command does
	Description() string

	// Execute executes the command and returns a tea.Cmd for Bubble Tea to process
	Execute(args string) tea.Cmd
}
