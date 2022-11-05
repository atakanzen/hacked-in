package main

import (
	"fmt"
	"os"

	"github.com/atakanzen/stacked-in/cli"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(cli.InitialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}
