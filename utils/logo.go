package utils

import (
	"fmt"
	"os"
	"os/exec"
	"github.com/charmbracelet/lipgloss"
)

// DisplayLogo prints the full CodeAid ASCII art logo
func DisplayLogo() {
	// Clear terminal before showing logo
	clearScreen()
	
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
		Render(logo)

	fmt.Println(styledLogo)
}

// clearScreen clears the terminal screen
func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
