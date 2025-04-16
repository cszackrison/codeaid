package utils

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

// DisplayLogo prints the CodeAid ASCII art logo
func DisplayLogo() {
	logo := `
   _____          _       _    _     _ 
  / ____|        | |     / \  (_)   | |
 | |     ___   __| | ___|  /\  _  __| |
 | |    / _ \ / _  |/ _ \ /  \| |/ _  |
 | |___| (_) | (_| |  __/ /\  \ | (_| |
  \_____\___/ \__,_|\___/_/  \_\_|\__,_|
                                       
`
	styledLogo := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("blue")).
		Render(logo)

	fmt.Println(styledLogo)
}
