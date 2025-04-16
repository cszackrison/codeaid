package utils

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

// ResponseMsg defines a custom message type for responses
type ResponseMsg string

// CancelMsg is a special message type returned when an operation is canceled
type CancelMsg struct{}

// Global client and environment setup to avoid repeated initialization
var (
	apiClient         *openai.Client
	apiKey            string
	clientInitMux     sync.Mutex
	clientInitialized bool
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

// FetchReply creates a tea.Cmd that fetches a reply with guaranteed completion
func FetchReply(prompt string) tea.Cmd {
	return func() tea.Msg {
		// This is the key change - we don't return nil from this function
		// Instead we use a timer to ensure we always return a result

		// Create result channel with buffer to avoid blocking
		resultChan := make(chan tea.Msg, 1)

		// Set up a failsafe timeout
		timer := time.NewTimer(15 * time.Second)
		defer timer.Stop()

		// Launch API call in goroutine
		go func() {
			client := initClient()

			// Create context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Make API request
			resp, err := client.CreateChatCompletion(
				ctx,
				openai.ChatCompletionRequest{
					Model:       "mistralai/mistral-small-3.1-24b-instruct:free",
					MaxTokens:   1024,
					Temperature: 0.7,
					Messages: []openai.ChatCompletionMessage{
						{
							Role:    openai.ChatMessageRoleUser,
							Content: prompt,
						},
					},
				},
			)

			// Process result or error
			if err != nil {
				resultChan <- ResponseMsg("Error: " + err.Error())
			} else if len(resp.Choices) == 0 {
				resultChan <- ResponseMsg("Error: No response received from API")
			} else {
				resultChan <- ResponseMsg(resp.Choices[0].Message.Content)
			}
		}()

		// Wait for either a result or timeout
		select {
		case result := <-resultChan:
			// Normal result path
			return result
		case <-timer.C:
			// Timeout path
			return ResponseMsg("Error: Request timed out after 15 seconds. Please try again.")
		}
	}
}
