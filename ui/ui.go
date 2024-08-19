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
	// Main page displays/changes. Picker page acts on it
	CurrentSource cfg.Source
	// TODO width needs to go here and be used by all pages
}

type Model struct {
	Common      *Common
	cfg         cfg.Config
	sourceList  SourcePickerModel
	releaseList ReleasePickerModel
	mainPage    FPageModel
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	// TODO: Proper resizing for all pages
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

func (m Model) View() string {
	switch m.Common.Page {
	case Main:
		return m.mainPage.View()

	case ReleasePicker:
		return m.releaseList.View()

	case SourcePicker:
		return m.sourceList.View()
	}

	return "\n FIXME"
}

// TODO this sounds like a dumb name and I need to handle errors here and/or main
func NewUI(cfg cfg.Config) Model {
	// shared state for all pages
	common := Common{Main, cfg.Sources[0]}

	// Create Source Picker Page (mostly wraps list Bubble)
	// Must copy to satisfy interface input(?)
	sourceItems := make([]list.Item, len(cfg.Sources))
	for i := range cfg.Sources {
		sourceItems[i] = cfg.Sources[i]
	}
	sl := NewSourcePicker(&common, sourceItems)

	// Create Main Page
	fp := NewFPModel(&common)

	// Create Release Picker Page
	rl := NewReleasePicker(&common)
	return Model{&common, cfg, sl, rl, fp}
}
