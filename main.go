
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	choice string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y":
			m.choice = "You chose YES!"
		case "n":
			m.choice = "You chose NO!"
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	md := "# Make a choice\n\nPress **y** for yes, **n** for no, or **q** to quit."
	if m.choice != "" {
		md += fmt.Sprintf("\n\n**Result**: %s", m.choice)
	}
	out, _ := glamour.Render(md, "dark")
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Foreground(lipgloss.Color("205"))
	return style.Render(out)
}

func main() {
	p := tea.NewProgram(model{})
	if err := p.Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
