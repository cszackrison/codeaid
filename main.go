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
	// Initialize the markdown renderer if not already done
	if m.markdownRenderer == nil {
		renderer, _ := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(100), // Adjust based on typical terminal width
		)
		m.markdownRenderer = renderer
	}
	return nil
}

// isControlChar checks if a string contains control characters
func isControlChar(s string) bool {
	if s == "" {
		return false
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
		switch msg.String() {
		case "ctrl+c", "esc":
			// Cancel current operation or exit
			if m.loading {
				// Stop loading and add a canceled message
				m.loading = false
				m.messages = append(m.messages, Message{Content: "(Canceled by user)", IsUser: true})
				// Force the UI to refresh
				return m, tea.Batch()
			}
			return m, tea.Quit

		case "enter":
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

			// Run loading animation and process user input (checking for commands)
			return m, tea.Batch(
				utils.TickAnimation(),
				utils.ProcessUserInput(userInput),
			)

		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}

		default:
			// Add typed character if not a control character
			if !isControlChar(msg.String()) {
				m.input += msg.String()
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
		ai:      lipgloss.NewStyle(),                                            // Default color
		error:   lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true), // Basic red from default palette
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
				// Check if content looks like markdown and render it if it does
				if containsMarkdown(msg.Content) && m.markdownRenderer != nil {
					renderedContent, err := m.markdownRenderer.Render(msg.Content)
					if err == nil {
						conversation.WriteString(renderedContent)
					} else {
						// Fall back to normal render if markdown parsing fails
						conversation.WriteString(styles.ai.Render(msg.Content))
					}
				} else {
					conversation.WriteString(styles.ai.Render(msg.Content))
				}
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

	// Render input prompt
	var prompt string
	inputText := "Enter your message: " + m.input
	if m.loading {
		prompt = styles.active.Render(inputText)
	} else {
		prompt = styles.input.Render(inputText)
	}

	// Combine all elements
	return fmt.Sprintf("%s\n\n%s%s", header, conversation.String(), prompt)
}

func main() {
	// Clear screen and display logo first
	utils.DisplayLogo()
	
	// Create and initialize the markdown renderer
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(100), // Adjust based on typical terminal width
	)
	if err != nil {
		fmt.Printf("Error initializing markdown renderer: %v\n", err)
		os.Exit(1)
	}

	// Create initial model
	initialModel := model{
		messages: []Message{},
		markdownRenderer: renderer,
	}

	// Create program with alternateScreen option for better performance
	p := tea.NewProgram(initialModel)

	// Run program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}