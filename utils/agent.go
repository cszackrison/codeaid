package utils

import (
	"codeaid/messages"
	"context"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

// Global variable to hold cancellation function
var currentCancelFunc context.CancelFunc

// Global client and environment setup to avoid repeated initialization
var (
	apiClient           *openai.Client
	apiKey              string
	clientInitMux       sync.Mutex
	clientInitialized   bool
	conversationMux     sync.Mutex
	conversationHistory []openai.ChatCompletionMessage
)

// Initialize API client once
func initClient() *openai.Client {
	clientInitMux.Lock()
	defer clientInitMux.Unlock()

	if !clientInitialized {
		_ = godotenv.Load()
		apiKey = os.Getenv("OPENROUTER_API_KEY")
		config := openai.DefaultConfig(apiKey)
		config.BaseURL = "https://openrouter.ai/api/v1"
		apiClient = openai.NewClientWithConfig(config)
		clientInitialized = true
	}

	return apiClient
}

// ClearHistory resets the conversation history
func ClearHistory() tea.Msg {
	conversationMux.Lock()
	defer conversationMux.Unlock()

	conversationHistory = nil
	return messages.ClearHistoryMsg{}
}

// AddMessageToHistory adds an assistant message to the conversation history
// This is called only after a response is successfully displayed and not canceled
func AddMessageToHistory(content string) {
	conversationMux.Lock()
	defer conversationMux.Unlock()

	assistantMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: content,
	}
	conversationHistory = append(conversationHistory, assistantMessage)
}

// ProcessUserInput handles user input and checks for commands
func ProcessUserInput(input string) tea.Cmd {
	trimmedInput := strings.TrimSpace(input)

	// Check if input is a command (starts with /)
	if len(trimmedInput) > 0 && trimmedInput[0] == '/' {
		// Find command handler
		return ExecuteCommand(trimmedInput)
	}

	// If not a command or command not found, process as a normal message
	return FetchReply(input)
}

// ExecuteCommand finds and executes a command
func ExecuteCommand(input string) tea.Cmd {
	// Extract command name (everything before the first space)
	cmdName := input

	if idx := strings.Index(input, " "); idx > 0 {
		cmdName = input[:idx]
		// We'll use the args when we implement command arguments
		// args = strings.TrimSpace(input[idx+1:])
	}

	// Handle built-in commands (we'll replace this when we implement commands properly)
	switch cmdName {
	case "/clear":
		return func() tea.Msg {
			return ClearHistory()
		}
	case "/help":
		return func() tea.Msg {
			helpText := "Available commands:\n" +
				"/clear - Clear conversation history\n" +
				"/help  - Show this help message\n" +
				"/exit  - Exit the application"
			return messages.ResponseMsg(helpText)
		}
	case "/exit":
		return tea.Quit
	}

	// Command not found
	return FetchReply(input)
}

// CancelCurrentRequest cancels any ongoing API request
func CancelCurrentRequest() {
	if currentCancelFunc != nil {
		currentCancelFunc() // Cancel the current context
	}
}

// FetchReply creates a tea.Cmd that fetches a reply with guaranteed completion
func FetchReply(prompt string) tea.Cmd {
	return func() tea.Msg {
		// Create result channel with buffer to avoid blocking
		resultChan := make(chan tea.Msg, 1)

		// Set up a failsafe timeout
		timer := time.NewTimer(15 * time.Second)
		defer timer.Stop()

		// Create a cancellable context and store its cancel function globally
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		// Store the cancel function so it can be called when user cancels
		currentCancelFunc = cancel

		// Launch API call in goroutine
		go func() {
			// Make sure to clean up
			defer cancel()
			defer func() { currentCancelFunc = nil }()

			client := initClient()

			// Add user message to history
			userMessage := openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			}

			// Update conversation history
			conversationMux.Lock()
			conversationHistory = append(conversationHistory, userMessage)
			messagesCopy := make([]openai.ChatCompletionMessage, len(conversationHistory))
			copy(messagesCopy, conversationHistory)
			conversationMux.Unlock()

			// Make API request with full conversation history
			resp, err := client.CreateChatCompletion(
				ctx,
				openai.ChatCompletionRequest{
					Model:       "mistralai/mistral-small-3.1-24b-instruct:free",
					MaxTokens:   1024,
					Temperature: 0.7,
					Messages:    messagesCopy,
				},
			)

			// Check if context was canceled before sending response
			select {
			case <-ctx.Done():
				// Context was canceled, don't send a response
				return
			default:
				// Context not canceled, proceed with normal response
				if err != nil {
					resultChan <- messages.ResponseMsg("Error: " + err.Error())
				} else if len(resp.Choices) == 0 {
					resultChan <- messages.ResponseMsg("Error: No response received from API")
				} else {
					// Send response - it will be added to history after displaying
					resultChan <- messages.ResponseMsg(resp.Choices[0].Message.Content)
				}
			}
		}()

		// Wait for either a result or timeout
		select {
		case result := <-resultChan:
			// Normal result path
			return result
		case <-timer.C:
			// Timeout path
			cancel() // Make sure to cancel the context on timeout
			return messages.ResponseMsg("Error: Request timed out after 15 seconds. Please try again.")
		case <-ctx.Done():
			// Context was canceled, return CancelMsg
			return messages.CancelMsg{}
		}
	}
}