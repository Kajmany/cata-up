package ui

import (
	_ "embed"

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
	ReleaseChannel ReleaseType
	help           help.Model
}

func NewFPModel(commonState *Common) FPageModel {
	return FPageModel{commonState, Experimental, help.New()}
}

func (m FPageModel) Init() tea.Cmd {
	return nil
}

func (m FPageModel) Update(msg tea.Msg) (FPageModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, fKeyMap.c.Quit):
			return m, tea.Quit
		case key.Matches(msg, fKeyMap.c.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, fKeyMap.Sources):
			m.Common.Page = SourcePicker
		case key.Matches(msg, fKeyMap.GetReleases):
			m.Common.Page = ReleasePicker
		case key.Matches(msg, fKeyMap.Channel):
			// I must be dumb because this toggle logic feels like magic
			m.ReleaseChannel = 1 - m.ReleaseChannel
		}
	}

	return m, cmd
}

func (m FPageModel) View() string {
	var output string
	srcLabel := lipgloss.JoinHorizontal(lipgloss.Center, underLineStyle.Render("Source:"), labelStyle.Render(m.Common.CurrentSource.Name))

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
	return output + m.help.View(fKeyMap)
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

type frontKeyMap struct {
	c           commonKeys
	Sources     key.Binding
	GetReleases key.Binding
	Channel     key.Binding
}

var fKeyMap = frontKeyMap{
	c: comKeys,
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
}

func (k frontKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.c.Help, k.c.Quit}
}

func (k frontKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Sources, k.GetReleases, k.Channel}, // first column
		{k.c.Help, k.c.Quit},                  // second column
	}
}
