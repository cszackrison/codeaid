package messages

// ResponseMsg defines a custom message type for responses
type ResponseMsg string

// CommandResponseMsg defines a message type for command responses (displayed but not sent to LLM)
type CommandResponseMsg string

// CancelMsg is a special message type returned when an operation is canceled
type CancelMsg struct{}

// ClearHistoryMsg is a message type to indicate history clearing
type ClearHistoryMsg struct{}

// TickMsg is sent when the animation needs to update
type TickMsg struct{}