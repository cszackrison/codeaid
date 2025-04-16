package utils

import (
	"fmt"
	"os"
	"os/exec"
	"github.com/charmbracelet/lipgloss"
)

func DisplayLogo() {
	clearScreen()
	
	logo := `
 ██████╗ ██████╗ ██████╗ ███████╗ █████╗ ██╗██████╗ 
██╔════╝██╔═══██╗██╔══██╗██╔════╝██╔══██╗██║██╔══██╗
██║     ██║   ██║██║  ██║█████╗  ███████║██║██║  ██║
██║     ██║   ██║██║  ██║██╔══╝  ██╔══██║██║██║  ██║
╚██████╗╚██████╔╝██████╔╝███████╗██║  ██║██║██████╔╝
 ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝╚═╝  ╚═╝╚═╝╚═════╝ 
`
	styledLogo := lipgloss.NewStyle().
		Bold(true).
		Render(logo)

	fmt.Println(styledLogo)
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
