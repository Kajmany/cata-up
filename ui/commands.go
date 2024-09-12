package ui

import (
	"github.com/Kajmany/cata-up/picker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/v63/github"
)

type (
	ErrMsg         struct{ error }
	NewReleasesMsg struct {
		releases []*github.RepositoryRelease
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

		return NewReleasesMsg{rels}
	}
}

/*
func (m ReleasePickerModel) cmdGetRelArtifact tea.Cmd {
  return func() tea.Msg {

  }
}
*/
