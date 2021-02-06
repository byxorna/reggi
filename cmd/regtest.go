package main

import (
	"fmt"
	"os"

	"github.com/byxorna/regtest/pkg/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	model, err := ui.New(os.Args[1:])
	if err != nil {
		fmt.Printf("Unable to initialize model: %v\n", err)
		os.Exit(1)
	}
	p := tea.NewProgram(*model)
	p.EnterAltScreen()
	err = p.Start()
	p.ExitAltScreen()
	if err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
