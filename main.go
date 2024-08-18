package main

import (
	"fmt"
	"os"

	"github.com/Kajmany/cata-up/cfg"
	"github.com/Kajmany/cata-up/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cfg, err := cfg.GetConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// Main Model w/ All Pages Attached
	m := ui.NewUI(cfg)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
