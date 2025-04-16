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

// Message represents a single chat message
type Message struct {
	Content string
	IsUser  bool
}

// Model represents the application state
type model struct {
	input           string
	cursorPosition  int
	messages        []Message
	loading         bool
	animationTick   int
	viewport        viewport
	markdownRenderer *glamour.TermRenderer
}

// Viewport manages the visible area of the chat
type viewport struct {
	width  int
	height int
}

func (m model) Init() tea.Cmd {
	// No markdown renderer initialization - disabled
	return nil
}

// isControlChar checks if a string contains control characters
func isControlChar(s string) bool {
	if s == "" {
		return false
	}

	// Check for special keys that should be treated as control characters
	specialKeys := []string{"up", "down", "left", "right", "home", "end", "delete", "del", "backspace", "enter", "tab", "esc"}
	for _, key := range specialKeys {
		if s == key {
			return true
		}
	}
	
	// Check for key combinations
	if strings.HasPrefix(s, "ctrl+") || strings.HasPrefix(s, "shift+") || strings.HasPrefix(s, "alt+") {
		return true
	}

	// Check for control sequences
	if strings.Contains(s, "\u0007") || strings.HasPrefix(s, "\u001b") || strings.HasPrefix(s, "\u001B") {
		return true
	}

	for _, r := range s {
		if unicode.IsControl(r) || r < 32 || (r >= 0x7F && r <= 0x9F) {
			return true
		}
	}

	return false
}

// containsMarkdown checks if content likely contains markdown formatting
func containsMarkdown(content string) bool {
	// Check for common markdown indicators
	markdownPatterns := []string{
		"```", // Code blocks
		"# ",  // Headers
		"## ", 
		"### ", 
		"* ",  // Lists
		"- ",
		"1. ", // Ordered lists
		"[",   // Links
		"![",  // Images
		"**",  // Bold
		"__",  // Underline
		"*",   // Italic
		"_",   // Italic
		">",   // Blockquotes
		"|",   // Tables
		"---", // Horizontal rule
	}

	for _, pattern := range markdownPatterns {
		if strings.Contains(content, pattern) {
			return true
		}
	}
	
	return false
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Use key types for all special keys for better reliability
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			// Cancel current operation or exit
			if m.loading {
				// Stop loading and add a canceled message
				m.loading = false
				m.messages = append(m.messages, Message{Content: "(Canceled by user)", IsUser: true})
				// Force the UI to refresh
				return m, tea.Batch()
			}
			return m, tea.Quit

		case tea.KeyEnter:
			if m.input == "" {
				return m, nil
			}

			// Cancel existing operation if needed
			if m.loading {
				m.loading = false
				m.messages = append(m.messages, Message{Content: "(Previous request canceled)", IsUser: true})
			}

			// Process new user input
			userInput := m.input
			m.messages = append(m.messages, Message{Content: userInput, IsUser: true})
			m.loading = true
			m.input = ""
			m.cursorPosition = 0

			// Run loading animation and process user input (checking for commands)
			return m, tea.Batch(
				utils.TickAnimation(),
				utils.ProcessUserInput(userInput),
			)

		case tea.KeyBackspace:
			if len(m.input) > 0 && m.cursorPosition > 0 {
				// Remove character before cursor
				m.input = m.input[:m.cursorPosition-1] + m.input[m.cursorPosition:]
				m.cursorPosition--
			}
			
		case tea.KeyDelete:
			if len(m.input) > 0 && m.cursorPosition < len(m.input) {
				// Remove character at cursor
				m.input = m.input[:m.cursorPosition] + m.input[m.cursorPosition+1:]
			}

		case tea.KeyLeft:
			if m.cursorPosition > 0 {
				m.cursorPosition--
			}
			
		case tea.KeyRight:
			if m.cursorPosition < len(m.input) {
				m.cursorPosition++
			}
			
		case tea.KeyHome:
			m.cursorPosition = 0
			
		case tea.KeyEnd:
			m.cursorPosition = len(m.input)
			
		case tea.KeyUp, tea.KeyDown:
			// Ignore up/down keys
			
		case tea.KeySpace:
			// Handle space key
			if m.cursorPosition == len(m.input) {
				m.input += " "
			} else {
				m.input = m.input[:m.cursorPosition] + " " + m.input[m.cursorPosition:]
			}
			m.cursorPosition++
			
		default:
			// For all other keys, check if they're text input
			if msg.Type == tea.KeyRunes {
				// Regular character input
				if m.cursorPosition == len(m.input) {
					m.input += string(msg.Runes)
				} else {
					m.input = m.input[:m.cursorPosition] + string(msg.Runes) + m.input[m.cursorPosition:]
				}
				m.cursorPosition++
			}
		}

	case utils.ResponseMsg:
		// Handle API response - check for error prefix
		content := string(msg)
		if strings.HasPrefix(content, "Error:") {
			// Handle error by showing it in the conversation with error styling
			m.messages = append(m.messages, Message{Content: content, IsUser: false})
			m.loading = false
			return m, nil
		}

		// Handle successful response
		m.messages = append(m.messages, Message{Content: content, IsUser: false})
		m.loading = false
		return m, nil

	case utils.ClearHistoryMsg:
		// Clear the chat history in the UI
		m.messages = []Message{
			{Content: "Chat history has been cleared.", IsUser: false},
		}
		m.loading = false
		return m, nil

	case utils.CancelMsg:
		// Handle explicit cancellation
		m.loading = false
		return m, nil

	case utils.TickMsg:
		// Update animation frame
		if m.loading {
			m.animationTick++
			return m, utils.TickAnimation()
		}
		return m, nil

	case tea.WindowSizeMsg:
		// Handle window resizing
		m.viewport.width = msg.Width
		m.viewport.height = msg.Height
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	// Define styles using default terminal colors where possible
	styles := struct {
		header  lipgloss.Style
		user    lipgloss.Style
		ai      lipgloss.Style
		error   lipgloss.Style
		loading lipgloss.Style
		input   lipgloss.Style
		active  lipgloss.Style
	}{
		header:  lipgloss.NewStyle().Bold(true),                                 // Bold, default color
		user:    lipgloss.NewStyle().Bold(true),                                 // Bold, default color
		ai:      lipgloss.NewStyle().Width(m.viewport.width - 2),                // Default color with width constraint
		error:   lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true).Width(m.viewport.width - 2), // Basic red with width constraint
		loading: lipgloss.NewStyle().Foreground(lipgloss.Color("3")),            // Basic yellow from default palette
		input:   lipgloss.NewStyle(),                                            // Default color
		active:  lipgloss.NewStyle().Underline(true),                            // Underlined instead of background color
	}

	// Render header
	header := styles.header.Render("CodeAid Chat")

	// Build message history using StringBuilder for better performance
	var conversation strings.Builder
	for _, msg := range m.messages {
		if msg.IsUser {
			conversation.WriteString(styles.user.Render("> " + msg.Content))
		} else {
			// Check if this is an error message
			if strings.HasPrefix(msg.Content, "Error:") {
				conversation.WriteString(styles.error.Render(msg.Content))
			} else {
				// No markdown rendering - display plain text with wrapping
				conversation.WriteString(styles.ai.Render(msg.Content))
			}
		}
		conversation.WriteString("\n\n")
	}

	// Add loading animation if active
	if m.loading {
		spinner := utils.GetLoadingAnimation(m.animationTick)
		conversation.WriteString(styles.loading.Render("Thinking " + spinner))
		conversation.WriteString("\n\n")
	}

	// Render input prompt with cursor
	var prompt string
	prefix := "Enter your message: "
	
	if m.loading {
		prompt = styles.active.Render(prefix + m.input)
	} else {
		// Display input with cursor indicator
		if m.cursorPosition == len(m.input) {
			// Cursor at the end
			prompt = styles.input.Render(prefix + m.input + "▎")
		} else {
			// Cursor in the middle - highlight the character at cursor position
			beforeCursor := m.input[:m.cursorPosition]
			atCursor := ""
			if m.cursorPosition < len(m.input) {
				atCursor = string(m.input[m.cursorPosition])
			}
			afterCursor := ""
			if m.cursorPosition+1 <= len(m.input) {
				afterCursor = m.input[m.cursorPosition+1:]
			}
			
			cursorStyle := lipgloss.NewStyle().Background(lipgloss.Color("7"))
			prompt = styles.input.Render(prefix + beforeCursor + cursorStyle.Render(atCursor) + afterCursor)
		}
	}

	// Combine all elements
	return fmt.Sprintf("%s\n\n%s%s", header, conversation.String(), prompt)
}

func main() {
	// Clear screen and display logo first
	utils.DisplayLogo()
	
	// Create initial model with default window size for proper text wrapping
	initialModel := model{
		messages:       []Message{},
		cursorPosition: 0,
		viewport: viewport{
			width:  80, // Default width, will be updated on first WindowSizeMsg
			height: 24, // Default height, will be updated on first WindowSizeMsg
		},
	}

	// Create program with alternateScreen option for better performance
	p := tea.NewProgram(initialModel)

	// Run program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}