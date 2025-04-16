package utils

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"os"
	"os/exec"
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
