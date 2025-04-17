package messages

// ResponseMsg defines a custom message type for responses
type ResponseMsg string

// CommandResponseMsg defines a message type for command responses (displayed but not sent to LLM)
type CommandResponseMsg string

// HelpMsg defines a special message type for the help command
type HelpMsg struct {
	Header   string
	Commands []CommandInfo
}

// CommandInfo holds information about a command
type CommandInfo struct {
	Name        string
	Description string
}

// CancelMsg is a special message type returned when an operation is canceled
type CancelMsg struct{}

// ClearHistoryMsg is a message type to indicate history clearing
type ClearHistoryMsg struct{}

// TickMsg is sent when the animation needs to update
type TickMsg struct{}

// ConfigMsg is a unified configuration message type with a type field
type ConfigMsg struct {
	Type        string      // "init", "prompt", "complete"
	PromptText  string      // For prompts and init messages
	Options     []string    // For selection prompts
	ConfigStep  string      // Current configuration step
	Config      interface{} // Configuration data
	FilePath    string      // For complete messages
	CurrentKey  string      // For init messages (masked key)
	CurrentModel string     // For init messages
}