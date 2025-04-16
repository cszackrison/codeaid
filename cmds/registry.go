package cmds

import (
	"strings"
)

// commandRegistry holds all registered commands
var commandRegistry []Command

// RegisterCommand adds a command to the registry
func RegisterCommand(cmd Command) {
	commandRegistry = append(commandRegistry, cmd)
}

// GetCommand returns a command by its name, or nil if not found
func GetCommand(name string) Command {
	for _, cmd := range commandRegistry {
		if cmd.Name() == name {
			return cmd
		}
	}
	return nil
}

// GetAllCommands returns all registered commands
func GetAllCommands() []Command {
	return commandRegistry
}

// GetCommandNames returns all command names
func GetCommandNames() []string {
	var names []string
	for _, cmd := range commandRegistry {
		names = append(names, cmd.Name())
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

	var matches []string
	for _, cmd := range commandRegistry {
		if strings.HasPrefix(cmd.Name(), prefix) {
			matches = append(matches, cmd.Name())
		}
	}
	return matches
}
