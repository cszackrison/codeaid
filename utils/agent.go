package utils

import (
	"context"
	"os"

	"github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

// FetchReply creates a tea.Cmd that fetches a reply from the OpenRouter API
func FetchReply(prompt string) tea.Cmd {
	return func() tea.Msg {
		_ = godotenv.Load()
		key := os.Getenv("OPENROUTER_API_KEY")
		config := openai.DefaultConfig(key)
		config.BaseURL = "https://openrouter.ai/api/v1"
		client := openai.NewClientWithConfig(config)

		resp, err := client.CreateChatCompletion(
			context.Background(),
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

		if err != nil {
			return "Error: " + err.Error()
		}

		if len(resp.Choices) == 0 {
			return "Error: No response received from API"
		}

		return resp.Choices[0].Message.Content
	}
}
