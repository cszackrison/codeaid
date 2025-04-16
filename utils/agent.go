package utils

import (
	"context"
	"os"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

// FetchReply creates a non-blocking tea.Cmd that fetches a reply from the OpenRouter API
func FetchReply(prompt string) tea.Cmd {
	return func() tea.Msg {
		// Create a channel to receive the API response
		responseChan := make(chan tea.Msg)

		// Execute API call in a goroutine to avoid blocking the UI
		go func() {
			_ = godotenv.Load()
			key := os.Getenv("OPENROUTER_API_KEY")
			config := openai.DefaultConfig(key)
			config.BaseURL = "https://openrouter.ai/api/v1"
			client := openai.NewClientWithConfig(config)

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			resp, err := client.CreateChatCompletion(
				ctx,
				openai.ChatCompletionRequest{
					Model:     "google/gemini-2.5-pro-exp-03-25:free",
					MaxTokens: 1024,
					Messages: []openai.ChatCompletionMessage{
						{
							Role:    openai.ChatMessageRoleUser,
							Content: prompt,
						},
					},
				},
			)

			var result tea.Msg
			if err != nil {
				result = "Error: " + err.Error()
			} else if len(resp.Choices) == 0 {
				result = "Error: No response received from API"
			} else {
				result = resp.Choices[0].Message.Content
			}

			responseChan <- result
		}()

		// Return from the Command function immediately, but with a pending channel read
		// This allows bubbletea to continue updating the UI while waiting for the API response
		return <-responseChan
	}
}
