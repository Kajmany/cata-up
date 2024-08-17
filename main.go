package main

import (
	"fmt"
	"os"

	"github.com/Kajmany/cata-up/cfg"
	"github.com/Kajmany/cata-up/ui"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type curPage int

const (
	Main curPage = iota
	SourcePicker
	ReleasePicker
)

type model struct {
	cfg         cfg.Config
	sourceList  list.Model
	releaseList ui.ReleasePickerModel
	mainPage    ui.FPageModel
	curPage     curPage
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.sourceList.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			// Use selected Source as new one, go back to Main
			if m.curPage == SourcePicker {
				i := m.sourceList.SelectedItem()
				m.mainPage.CurrentSource = i.(cfg.Source)
				m.curPage = Main
			}
			return m, nil

		case "s":
			// Open source picker from main page
			if m.curPage == Main {
				m.curPage = SourcePicker
				return m, nil
			}

		case "c":
			// Change relase channel from main page
			if m.curPage == Main {
				if m.mainPage.ReleaseChannel == ui.Experimental {
					m.mainPage.ReleaseChannel = ui.Stable
				} else {
					m.mainPage.ReleaseChannel = ui.Experimental
				}
				return m, nil
			}

		case "g":
			// Change to release picker from main
			if m.curPage == Main {
				m.curPage = ReleasePicker
			}

		case "esc", "backspace":
			// Go back to main page from pickers screens
			if m.curPage == ReleasePicker || m.curPage == SourcePicker {
				m.curPage = Main
			}
		}

		var cmd tea.Cmd
		m.sourceList, cmd = m.sourceList.Update(msg)
		return m, cmd

	}
	return m, nil // TODO fall through
}

func (m model) View() string {
	switch m.curPage {
	case Main:
		return m.mainPage.View()

	case ReleasePicker:
		return m.releaseList.View()

	case SourcePicker:
		return "\n" + m.sourceList.View()
	}

	return "\n FIXME"
}

func main() {
	cfg, err := cfg.GetConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// Create Source Picker List Bubbble
	// Must copy to satisfy interface input
	sourceItems := make([]list.Item, len(cfg.Sources))
	for i := range cfg.Sources {
		sourceItems[i] = cfg.Sources[i]
	}
	sl := list.New(sourceItems, list.NewDefaultDelegate(), 20, 14)
	sl.Title = "Sources"

	// Create Main Page
	fp := ui.NewFPModel(cfg.Sources[0])

	// Create Release Picker Page
	rl := ui.NewReleasePicker()

	// Main Model w/ All Pages Attached
	m := model{cfg, sl, rl, fp, Main}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
