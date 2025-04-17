package cmds

import (
	"codeaid/utils"
	tea "github.com/charmbracelet/bubbletea"
)

// ClearCommand clears the conversation history
type ClearCommand struct{}

// Name returns the command name
func (c ClearCommand) Name() string {
	return "/clear"
}

// Description returns the command description
func (c ClearCommand) Description() string {
	return "Clear conversation history"
}

// Execute executes the command
func (c ClearCommand) Execute(args string) tea.Cmd {
	return func() tea.Msg {
		return utils.ClearHistory()
	}
}

