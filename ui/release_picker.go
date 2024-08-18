package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ReleasePickerModel struct {
	// Will need to swap between fetched release types on request
	ExperimentalReleases []list.Item
	StableReleases       []list.Item
	list                 list.Model
}

// TODO is this nonsensical compared to other init workflows
func NewReleasePicker() ReleasePickerModel {
	var rp ReleasePickerModel
	rp.initReleasePicker()
	return rp
}

func (m ReleasePickerModel) Init() tea.Cmd {
	return nil
}

func (m *ReleasePickerModel) initReleasePicker() {
	// TODO handle width/height
	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), 40, 20)
	m.list.Title = "Available Releases"
}

func (m ReleasePickerModel) Update(msg tea.Msg) (ReleasePickerModel, tea.Cmd) {
	var cmd tea.Cmd
	return m, cmd
}

func (m ReleasePickerModel) View() string {
	return m.list.View()
}
