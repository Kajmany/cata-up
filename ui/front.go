package ui

import (
	"github.com/Kajmany/cata-up/cfg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	_ "embed"
)

type ReleaseType int

const (
	Experimental ReleaseType = iota
	Stable
)

type FPageModel struct {
	CurrentSource  cfg.Source
	ReleaseChannel ReleaseType
}

func NewFPModel(defaultSource cfg.Source) FPageModel {
	return FPageModel{defaultSource, Experimental}
}

func (m FPageModel) Init() tea.Cmd {
	return nil
}

func (m FPageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "s":
			// i, ok = m.sourceList.
			return m, nil
		}
	}

	var cmd tea.Cmd
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
	return output
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
