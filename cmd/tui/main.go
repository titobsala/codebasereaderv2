package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tito-sala/codebasereaderv2/internal/tui/core"
)

func main() {
	// Create the main TUI model
	model := core.NewMainModel()

	// Create the Bubble Tea program
	p := tea.NewProgram(model, tea.WithAltScreen())

	// Run the program
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Thanks for using CodebaseReader v2!")
}
