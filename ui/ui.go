package ui

// TODO this file is mostly junk ripped out of Main
import (
	"github.com/Kajmany/cata-up/cfg"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type CurPage int

const (
	Main CurPage = iota
	SourcePicker
	ReleasePicker
)

type Common struct {
	Page CurPage
	// TODO width needs to go here and be used by all pages
}

type Model struct {
	Common      *Common
	cfg         cfg.Config
	sourceList  list.Model // TODO make this a custom page just for consistency?
	releaseList ReleasePickerModel
	mainPage    FPageModel
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.sourceList.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		// TODO move the rest of this to page-specific update functions
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			// Use selected Source as new one, go back to Main
			if m.Common.Page == SourcePicker {
				i := m.sourceList.SelectedItem()
				m.mainPage.CurrentSource = i.(cfg.Source)
				m.Common.Page = Main
			}
			return m, nil

		case "esc", "backspace":
			// Go back to main page from pickers screens
			if m.Common.Page == ReleasePicker || m.Common.Page == SourcePicker {
				m.Common.Page = Main
			}
		}
		switch m.Common.Page {
		case Main:
			m.mainPage, cmd = m.mainPage.Update(msg)
			cmds = append(cmds, cmd)
		case SourcePicker:
			m.sourceList, cmd = m.sourceList.Update(msg)
			cmds = append(cmds, cmd)
		case ReleasePicker:
			m.releaseList, cmd = m.releaseList.Update(msg)
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)

	}
	return m, nil // TODO fall through
}

func (m Model) View() string {
	switch m.Common.Page {
	case Main:
		return m.mainPage.View()

	case ReleasePicker:
		return m.releaseList.View()

	case SourcePicker:
		return "\n" + m.sourceList.View()
	}

	return "\n FIXME"
}

// TODO this sounds like a dumb name and I need to handle errors here and/or main
func NewUI(cfg cfg.Config) Model {
	// shared state for all pages
	common := Common{Main}

	// Create Source Picker List Bubbble
	// Must copy to satisfy interface input
	sourceItems := make([]list.Item, len(cfg.Sources))
	for i := range cfg.Sources {
		sourceItems[i] = cfg.Sources[i]
	}
	sl := list.New(sourceItems, list.NewDefaultDelegate(), 20, 14)
	sl.Title = "Sources"

	// Create Main Page
	fp := NewFPModel(&common, cfg.Sources[0])

	// Create Release Picker Page
	rl := NewReleasePicker()
	return Model{&common, cfg, sl, rl, fp}
}
