package ui

import (
	_ "embed"

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
	c              *Common
	ReleaseChannel ReleaseType

	KeyMap frontKeyMap
}

func NewFPModel(commonState *Common) FPageModel {
	return FPageModel{commonState, Experimental, fKeyMap}
}

func (m FPageModel) Init() tea.Cmd {
	return nil
}

func (m FPageModel) Update(msg tea.Msg) (FPageModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, fKeyMap.Sources):
			m.c.Page = SourcePicker
		case key.Matches(msg, fKeyMap.GetReleases):
			m.c.Page = ReleasePicker
		case key.Matches(msg, fKeyMap.Channel):
			// I must be dumb because this toggle logic feels like magic
			m.ReleaseChannel = 1 - m.ReleaseChannel
		}
	}

	return m, cmd
}

func (m FPageModel) View() string {
	srcLabel := lipgloss.JoinHorizontal(lipgloss.Center, underLineStyle.Render("Source:"), labelStyle.Render(m.c.CurrentSource.Name))

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
	content := lipgloss.JoinVertical(lipgloss.Center, title, srcLabel, channelRadio)
	return lipgloss.JoinVertical(lipgloss.Center, content, m.c.help.View(m))
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

func (k FPageModel) ShortHelp() []key.Binding {
	return []key.Binding{k.KeyMap.c.Help, k.KeyMap.c.Quit}
}

func (k FPageModel) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.KeyMap.Sources, k.KeyMap.GetReleases, k.KeyMap.Channel}, // first column
		{k.KeyMap.c.Help, k.KeyMap.c.Quit},                         // second column
	}
}
