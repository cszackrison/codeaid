package cmds

import (
	"codeaid/utils"
	"strings"
)

// commandRegistry holds all registered commands
var commandRegistry = make(map[string]Command)

// init function registers all commands
func init() {
	// Register all commands here
	RegisterCommand(ClearCommand{})
	RegisterCommand(ExitCommand{})
	RegisterCommand(HelpCommand{})
	RegisterCommand(ConfigCommand{})
}

// RegisterCommand adds a command to the registry
func RegisterCommand(cmd Command) {
	commandRegistry[cmd.Name()] = cmd
}

// GetCommand returns a command by its name, or nil if not found
func GetCommand(name string) utils.CommandExecutor {
	cmd, ok := commandRegistry[name]
	if !ok {
		return nil
	}
	return cmd
}

// GetAllCommands returns all registered commands
func GetAllCommands() []Command {
	cmds := make([]Command, 0, len(commandRegistry))
	for _, cmd := range commandRegistry {
		cmds = append(cmds, cmd)
	}
	return cmds
}

// GetCommandNames returns all command names
func GetCommandNames() []string {
	names := make([]string, 0, len(commandRegistry))
	for name := range commandRegistry {
		names = append(names, name)
	}
	return names
}

// FindMatchingCommands returns all commands that start with the given prefix
func FindMatchingCommands(prefix string) []string {
	if prefix == "" {
		return nil
	}

	// Ensure prefix starts with /
	if !strings.HasPrefix(prefix, "/") {
		return nil
	}

	matches := []string{}
	for name := range commandRegistry {
		if strings.HasPrefix(name, prefix) {
			matches = append(matches, name)
		}
	}
	return matches
}

// CommandRegistry implements the utils.CommandHandler interface
type CommandRegistry struct{}

// GetCommand returns a command by its name
func (cr CommandRegistry) GetCommand(name string) utils.CommandExecutor {
	return GetCommand(name)
}