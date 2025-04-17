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

// ConfigInitMsg is sent to start the config setup process
type ConfigInitMsg struct {
	CurrentAPIKey string
	CurrentModel  string 
}

// ConfigPromptMsg is sent to prompt for a specific config value
type ConfigPromptMsg struct {
	PromptText string
	Options    []string
	ConfigStep string
	Config     interface{}
}

// ConfigInputMsg is sent when user provides input for a config prompt
type ConfigInputMsg struct {
	Input      string
	ConfigStep string
	Config     interface{}
}

// ConfigCompleteMsg is sent when config setup is complete
type ConfigCompleteMsg struct {
	FilePath string
}