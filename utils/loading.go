package utils

import (
	"codeaid/messages"
	"time"

	"github.com/charmbracelet/bubbletea"
)

// Animation constants
const (
	// Using a shorter animation interval for smoother animation
	animationInterval = time.Millisecond * 60
)

// Simple animation characters that work with any terminal
var animChars = []string{"⠋", "⠙", "⠚", "⠞", "⠖", "⠦", "⠴", "⠲", "⠳", "⠓"}

// TickAnimation sends a tick message after a short delay
func TickAnimation() tea.Cmd {
	return tea.Tick(animationInterval, func(_ time.Time) tea.Msg {
		return messages.TickMsg{}
	})
}

// GetLoadingAnimation returns the current loading animation character
func GetLoadingAnimation(tick int) string {
	return animChars[tick%len(animChars)]
}
