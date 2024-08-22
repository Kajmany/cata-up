package ui

// TODO this file is mostly junk ripped out of Main
import (
	"github.com/Kajmany/cata-up/cfg"
	"github.com/Kajmany/cata-up/log"
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
	logger log.BufferedLogger
}

type Model struct {
	Common      *Common
	cfg         cfg.Config
	sourceList  SourcePickerModel
	releaseList ReleasePickerModel
	mainPage    FPageModel
}

func (m Model) Init() tea.Cmd {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)
	// Need to remember to add init for all new pages
	cmd = m.mainPage.Init()
	cmds = append(cmds, cmd)
	cmd = m.sourceList.Init()
	cmds = append(cmds, cmd)
	m.Common.logger.L.Debug("Preparing to call RL init")
	cmd = m.releaseList.Init()
	cmds = append(cmds, cmd)
	m.Common.logger.L.Debug("Returning a batch of commands", "number", len(cmds))
	return tea.Batch(cmds...)
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
	// Some async messages may need to go to pages that aren't active
	// TODO: can this be done better?
	if msg, ok := msg.(NewReleasesMsg); ok {
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
func NewUI(cfg cfg.Config, logger log.BufferedLogger) Model {
	// shared state for all pages
	common := Common{Main, cfg.Sources[0], logger}

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
	common.logger.L.Info("Welcome to cata-up!")
	return Model{&common, cfg, sl, rl, fp}
}
