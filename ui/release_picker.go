package ui

import (
	"github.com/Kajmany/cata-up/picker"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/v63/github"
)

type ReleasePickerModel struct {
	common *Common
	// Will need to swap between fetched release types on request
	ExperimentalReleases []list.Item
	StableReleases       []list.Item // TODO not handled
	client               *github.Client
	list                 list.Model
}

func (m ReleasePickerModel) Init() tea.Cmd {
	return m.cmdGetReleases()
}

func NewReleasePicker(common *Common) ReleasePickerModel {
	var m ReleasePickerModel
	m.common = common
	// TODO handle width/height
	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), 40, 20)
	m.list.Title = "Available Releases"
	m.client = picker.GetClient()
	return m
}

func (m ReleasePickerModel) Update(msg tea.Msg) (ReleasePickerModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case NewReleasesMsg:
		m.common.logger.L.Info("recieved releases", "number", len(msg.releases))
		m.ExperimentalReleases = append(m.ExperimentalReleases, msg.releases...)
		cmd = m.list.SetItems(m.ExperimentalReleases)

	case ErrMsg:
		m.common.logger.L.Error("update got error", "err", msg)
		return m, tea.Quit

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, comKeys.Quit):
			return m, tea.Quit

		case key.Matches(msg, comKeys.Back):
			m.common.Page = Main

		default:
			// FIXME: This doesn't work.
			m.list.Update(msg) // Pass to list underlying
		}
	}
	return m, cmd
}

func (m ReleasePickerModel) View() string {
	return m.list.View()
}

type (
	ErrMsg         struct{ error }
	NewReleasesMsg struct {
		releases []list.Item
	}
)

type Release struct{ *github.RepositoryRelease }

// Bubbletea list.Item methods
func (r Release) FilterValue() string {
	return *r.TagName
}

func (r Release) Title() string {
	return *r.TagName
}

func (r Release) Description() string {
	return r.CreatedAt.String()
}

// End Bubbletea methods

func (m ReleasePickerModel) cmdGetReleases() tea.Cmd {
	return func() tea.Msg {
		// TODO needs to be capable of getting more than 1st load
		rels, err := picker.GetRecentReleases(m.client, m.common.CurrentSource, 0, 10)
		if err != nil {
			// TODO make this error more useful for program
			return ErrMsg{err}
		}
		// FIXME is this dance of wrapping and then converting necessary? feels dumb
		relItems := make([]list.Item, len(rels))
		for i := range rels {
			relItems[i] = Release{rels[i]}
		}
		return NewReleasesMsg{relItems}
	}
}

/* Using common until page-specific bindings required
type releaseKeyMap struct {
	c commonKeys
}
*/
