package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"codeaid/messages"
	"codeaid/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// Message represents a single chat message
type Message struct {
	Content     string
	IsUser      bool
	IsCommand   bool
}

// Model represents the application state
type model struct {
	input            string
	cursorPosition   int
	messages         []Message
	loading          bool
	animationTick    int
	viewport         viewport
	markdownRenderer *glamour.TermRenderer
	hints            []string
	selectedHint     int
	showHints        bool
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
		"* ", // Lists
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
				// Stop loading without adding a response to the conversation
				m.loading = false
				// Cancel the actual API request first
				utils.CancelCurrentRequest()
				return m, nil
			}
			return m, tea.Quit

		case tea.KeyEnter:
			// If hints are shown and a hint is selected, use it instead
			if m.showHints && len(m.hints) > 0 && m.selectedHint >= 0 && m.selectedHint < len(m.hints) {
				m.input = m.hints[m.selectedHint]
				m.cursorPosition = len(m.input)
				m.showHints = false
				return m, nil
			}

			if m.input == "" {
				return m, nil
			}

			// Cancel existing operation if needed
			if m.loading {
				// Cancel the actual API request first
				utils.CancelCurrentRequest()
				// Set loading to false
				m.loading = false
				// Process the new input right away
				return m, utils.ProcessUserInput(m.input)
			}

			// Process new user input
			userInput := m.input
			m.messages = append(m.messages, Message{Content: userInput, IsUser: true})
			m.loading = true
			m.input = ""
			m.cursorPosition = 0
			m.showHints = false

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

				// Update hints based on new input
				m.hints = getCommandHints(m.input)
				if len(m.hints) > 0 {
					m.showHints = true
					m.selectedHint = 0
				} else {
					m.showHints = false
				}
			}

		case tea.KeyDelete:
			if len(m.input) > 0 && m.cursorPosition < len(m.input) {
				// Remove character at cursor
				m.input = m.input[:m.cursorPosition] + m.input[m.cursorPosition+1:]

				// Update hints based on new input
				m.hints = getCommandHints(m.input)
				if len(m.hints) > 0 {
					m.showHints = true
					m.selectedHint = 0
				} else {
					m.showHints = false
				}
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

		case tea.KeyUp:
			// Handle up key for hint navigation
			if m.showHints && len(m.hints) > 0 {
				if m.selectedHint > 0 {
					m.selectedHint--
				} else {
					// Wrap around to the last hint
					m.selectedHint = len(m.hints) - 1
				}
			}

		case tea.KeyDown:
			// Handle down key for hint navigation
			if m.showHints && len(m.hints) > 0 {
				if m.selectedHint < len(m.hints)-1 {
					m.selectedHint++
				} else {
					// Wrap around to the first hint
					m.selectedHint = 0
				}
			}

		case tea.KeySpace:
			// Handle space key
			if m.cursorPosition == len(m.input) {
				m.input += " "
			} else {
				m.input = m.input[:m.cursorPosition] + " " + m.input[m.cursorPosition:]
			}
			m.cursorPosition++

		case tea.KeyTab:
			// Handle tab key for autocomplete
			if m.showHints && len(m.hints) > 0 && m.selectedHint >= 0 && m.selectedHint < len(m.hints) {
				// Autocomplete with the selected hint
				m.input = m.hints[m.selectedHint]
				m.cursorPosition = len(m.input)
				m.showHints = false
			}

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

				// Update hints based on new input
				m.hints = getCommandHints(m.input)
				if len(m.hints) > 0 {
					m.showHints = true
					m.selectedHint = 0
				} else {
					m.showHints = false
				}
			}
		}

	case messages.ResponseMsg:
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
		// Now that we're displaying the response, update the conversation history
		return m, func() tea.Msg {
			// Add message to conversation history after it's displayed
			utils.AddMessageToHistory(content)
			return nil
		}

	case messages.CommandResponseMsg:
		// Handle command response (display only, not added to conversation history)
		content := string(msg)
		m.messages = append(m.messages, Message{Content: content, IsUser: false, IsCommand: true})
		m.loading = false
		return m, nil
		
	case messages.HelpMsg:
		// Handle help message with custom styling
		helpMsg := msg
		var sb strings.Builder
		
		// Style the header using styles from the View function
		headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("3"))
		sb.WriteString(headerStyle.Render(helpMsg.Header) + "\n")
		
		// Style for command names
		cmdStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4"))
		// Style for descriptions
		descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
		
		for _, cmd := range helpMsg.Commands {
			sb.WriteString(cmdStyle.Render(cmd.Name))
			sb.WriteString(" - ")
			sb.WriteString(descStyle.Render(cmd.Description))
			sb.WriteString("\n")
		}
		
		// Add the styled help content to messages
		m.messages = append(m.messages, Message{Content: sb.String(), IsUser: false, IsCommand: true})
		m.loading = false
		return m, nil

	case messages.ClearHistoryMsg:
		// Clear the chat history in the UI
		m.messages = []Message{
		}
		m.loading = false
		return m, nil

	case messages.CancelMsg:
		// The request was canceled via context cancellation
		// Just set loading to false without adding any message to UI or history
		m.loading = false
		return m, nil

	case messages.TickMsg:
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
		header       lipgloss.Style
		user         lipgloss.Style
		ai           lipgloss.Style
		error        lipgloss.Style
		loading      lipgloss.Style
		input        lipgloss.Style
		active       lipgloss.Style
		hint         lipgloss.Style
		hintSelected lipgloss.Style
		command      lipgloss.Style
	}{
		header:       lipgloss.NewStyle().Bold(true).Width(m.viewport.width - 2),
		user:         lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Width(m.viewport.width - 2),
		ai:           lipgloss.NewStyle().Width(m.viewport.width - 2),
		error:        lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true).Width(m.viewport.width - 2),
		loading:      lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		input:        lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("8")).MarginBottom(1).Width(m.viewport.width - 2),
		active:       lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("0")).Width(m.viewport.width - 2),
		hint:         lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Width(m.viewport.width - 2),
		hintSelected: lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Width(m.viewport.width - 2),
		command:      lipgloss.NewStyle().Foreground(lipgloss.Color("8")).PaddingLeft(4).Width(m.viewport.width - 2),
	}

	// Build message history using StringBuilder for better performance
	var conversation strings.Builder
	for _, msg := range m.messages {
		if msg.IsUser {
			conversation.WriteString(styles.user.Render("> " + msg.Content))
		} else if msg.IsCommand {
			// Don't apply any styling for command messages as they're already styled
			conversation.WriteString(styles.command.Render(msg.Content))
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
	prefix := "> "

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

	// Add hints if available
	var hintsDisplay string
	if m.showHints && len(m.hints) > 0 {
		var hintsBuilder strings.Builder
		hintsBuilder.WriteString("\n")

		for i, hint := range m.hints {
			if i == m.selectedHint {
				// Highlight the selected hint
				hintsBuilder.WriteString(styles.hintSelected.Render(" " + hint + " "))
			} else {
				hintsBuilder.WriteString(styles.hint.Render(" " + hint + " "))
			}
			hintsBuilder.WriteString("\n") // Add newline for vertical display
		}
		hintsDisplay = hintsBuilder.String()
	}

	// Combine all elements
	if m.showHints && len(m.hints) > 0 {
		return fmt.Sprintf("%s\n\n%s%s", conversation.String(), prompt, hintsDisplay)
	} else {
		return fmt.Sprintf("%s\n\n%s", conversation.String(), prompt)
	}
}

// getCommandHints returns a list of command hints that match the current input
func getCommandHints(input string) []string {
	// If input is empty or doesn't start with '/', return no hints
	if len(input) == 0 || input[0] != '/' {
		return nil
	}

	// Available commands
	commands := []string{
		"/clear",
		"/help",
		"/exit",
	}

	// Find matching commands
	var matches []string
	for _, cmd := range commands {
		if strings.HasPrefix(cmd, input) {
			matches = append(matches, cmd)
		}
	}

	return matches
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
		hints:        []string{},
		selectedHint: -1,
		showHints:    false,
	}

	// Create program with alternateScreen option for better performance
	p := tea.NewProgram(initialModel)

	// Run program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
