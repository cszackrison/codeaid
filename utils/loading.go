package utils

import (
	"time"

	"github.com/charmbracelet/bubbletea"
)

// TickMsg is sent when the animation needs to update
type TickMsg struct{}

// TickAnimation sends a tick message after a short delay
func TickAnimation() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return TickMsg{}
	})
}

// GetLoadingAnimation returns the current loading animation character
func GetLoadingAnimation(tick int) string {
	animChars := []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
	return animChars[tick%len(animChars)]
}
