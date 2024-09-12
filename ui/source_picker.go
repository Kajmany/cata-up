package ui

import (
	"github.com/Kajmany/cata-up/cfg"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Very little to add to built-in list class. Just need a way to go back
// So this is mostly wrapped for consistency w/ other pages

type SourcePickerModel struct {
	common *Common
	list   list.Model
}

func NewSourcePicker(common *Common, items []list.Item) SourcePickerModel {
	var m SourcePickerModel
	m.common = common
	m.list = list.New(items, list.NewDefaultDelegate(), 20, 40)
	m.list.SetShowHelp(false)
	m.list.Title = "Sources"
	return m
}

func (m SourcePickerModel) Init() tea.Cmd {
	return nil
}

func (m SourcePickerModel) Update(msg tea.Msg) (SourcePickerModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetHeight(msg.Height)
		m.list.SetWidth(msg.Width)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, comKeys.Back):
			if !m.list.SettingFilter() {
				m.common.Page = Main
				return m, nil
			}
		case key.Matches(msg, comKeys.Select):
			i := m.list.SelectedItem()
			m.common.CurrentSource = i.(cfg.Source) // TODO: dangerous mutation?
			m.common.Page = Main
			return m, nil
		}

		// And pass everything to the list, too.
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}
	// What type of message could this be? Doesn't matter.
	return m, cmd
}

func (m SourcePickerModel) View() string {
	heightBudget := m.common.height
	helpText := m.common.help.View(m)
	heightBudget -= lipgloss.Height(helpText)
	if heightBudget < 10 {
		// Totally arbitrary set point. Mostly concerned about negative.
		m.common.logger.L.Warn("Height available after help text < 10! Expect malformed window.",
			"budget", heightBudget)
	}

	// Width doesn't need to be calculated each view. It's set on update
	m.list.SetHeight(heightBudget)

	return lipgloss.JoinVertical(lipgloss.Center, m.list.View(), helpText)
}
