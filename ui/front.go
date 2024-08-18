package ui

import (
	_ "embed"

	"github.com/Kajmany/cata-up/cfg"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ReleaseType int

const (
	Experimental ReleaseType = iota
	Stable
)

type FPageModel struct {
	Common         *Common
	CurrentSource  cfg.Source
	ReleaseChannel ReleaseType
	help           help.Model
}

func NewFPModel(commonState *Common, defaultSource cfg.Source) FPageModel {
	return FPageModel{commonState, defaultSource, Experimental, help.New()}
}

func (m FPageModel) Init() tea.Cmd {
	return nil
}

func (m FPageModel) Update(msg tea.Msg) (FPageModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymap.Quit):
			return m, tea.Quit
		case key.Matches(msg, keymap.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, keymap.Sources):
			m.Common.Page = SourcePicker
		case key.Matches(msg, keymap.GetReleases):
			m.Common.Page = ReleasePicker
		case key.Matches(msg, keymap.Channel):
			// I must be dumb because this toggle logic feels like magic
			m.ReleaseChannel = 1 - m.ReleaseChannel
		}
	}

	return m, cmd
}

func (m FPageModel) View() string {
	var output string
	srcLabel := lipgloss.JoinHorizontal(lipgloss.Center, underLineStyle.Render("Source:"), labelStyle.Render(m.CurrentSource.Name))

	var (
		expStatus  string
		stabStatus string
	)
	// This feels so dumb
	if m.ReleaseChannel == Experimental {
		expStatus = "[X]"
		stabStatus = "[ ]"
	} else {
		expStatus = "[ ]"
		stabStatus = "[X]"
	}
	channelRadio := underLineStyle.Render("Channel") + "\n" + labelStyle.Render("Experimental: "+expStatus+"\nStable: "+stabStatus)
	output = lipgloss.JoinVertical(lipgloss.Center, title, srcLabel, channelRadio)
	return output + m.help.View(keymap)
}

var labelStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	Padding(0).
	BorderTop(true).
	BorderLeft(true).
	BorderRight(true).
	BorderBottom(true)

var underLineStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	Padding(0).
	BorderTop(false).
	BorderLeft(false).
	BorderRight(false).
	BorderBottom(true)

// FIGlet standard. What else?
//
//go:embed title.txt
var title string

type keyMap struct {
	Sources     key.Binding
	GetReleases key.Binding
	Channel     key.Binding
	Back        key.Binding
	Help        key.Binding
	Quit        key.Binding
}

var keymap = keyMap{
	Sources: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "source Selection"),
	),
	GetReleases: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "get Releases"),
	),
	Channel: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "change Release Channel"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc", "backspace"),
		key.WithHelp("esc/backspace", "back to Main"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Sources, k.GetReleases, k.Channel}, // first column
		{k.Help, k.Quit},                      // second column
	}
}
