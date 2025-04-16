package utils

import (
	"time"

	"github.com/charmbracelet/bubbletea"
)

// TickMsg is sent when the animation needs to update
type TickMsg struct{}

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
		return TickMsg{}
	})
}

// GetLoadingAnimation returns the current loading animation character
func GetLoadingAnimation(tick int) string {
	return animChars[tick%len(animChars)]
}
