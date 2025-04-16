package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"codeaid/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	input         string
	messages      []string
	loading       bool
	done          bool
	animationTick int
}

func (m model) Init() tea.Cmd {
	return nil
}

// isControlChar checks if a string is a control character
func isControlChar(s string) bool {
	if s == "" {
		return false
	}

	// Check for ANSI escape sequences
	if strings.HasPrefix(s, "\u001b") || strings.HasPrefix(s, "\u001B") {
		return true
	}

	// Check if it's a control character
	for _, r := range s {
		if unicode.IsControl(r) {
			return true
		}
	}

	return false
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		// Handle exit keys immediately
		if key == "esc" || key == "ctrl+c" {
			return m, tea.Quit
		}

		switch key {
		case "enter":
			if m.input != "" {
				userInput := m.input
				m.messages = append(m.messages, "> "+userInput)
				m.loading = true
				m.input = ""
				return m, tea.Batch(fetchReply(userInput), utils.TickAnimation())
			}
		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		default:
			// Filter out control characters that might be causing issues
			if !isControlChar(msg.String()) {
				m.input += msg.String()
			}
		}
	case string:
		m.messages = append(m.messages, msg)
		m.loading = false
		m.input = ""
		return m, nil
	case utils.TickMsg:
		if m.loading {
			m.animationTick++
			return m, utils.TickAnimation()
		}
	}
	return m, nil
}

func (m model) View() string {
	header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63")).Render("🤖 OpenRouter Chat")

	// Build the conversation history
	var conversation string
	if len(m.messages) > 0 {
		md, _ := glamour.NewTermRenderer(glamour.WithAutoStyle())
		for _, msg := range m.messages {
			if msg[:2] == "> " {
				// User message
				userMsg := lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Render(msg)
				conversation += userMsg + "\n"
			} else {
				// AI response
				formatted, _ := md.Render(msg)
				conversation += formatted + "\n"
			}
		}
	}

	// Loading animation
	if m.loading {
		currentChar := utils.GetLoadingAnimation(m.animationTick)
		loadingMsg := lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Render("Thinking " + currentChar)
		return header + "\n" + conversation + loadingMsg
	}

	// Input prompt
	prompt := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("Enter your message: " + m.input)
	return header + "\n" + conversation + prompt
}

func fetchReply(prompt string) tea.Cmd {
	return utils.FetchReply(prompt)
}

func main() {
	// Display the CodeAid logo
	utils.DisplayLogo()

	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
