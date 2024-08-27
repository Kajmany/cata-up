package ui

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/Kajmany/cata-up/common"
	"github.com/Kajmany/cata-up/picker"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/go-github/v63/github"
)

type examineMode int

const (
	showRelease examineMode = iota
	showChangelog
)

type ReleasePickerModel struct {
	common *Common
	// Will need to swap between fetched release types on request
	ExperimentalReleases []*github.RepositoryRelease
	StableReleases       []list.Item // TODO not handled
	client               *github.Client
	list                 list.Model
	port                 viewport.Model
	examiner             examineMode
	changelog            string
}

func (m ReleasePickerModel) Init() tea.Cmd {
	return m.cmdGetReleases()
}

func NewReleasePicker(common *Common) ReleasePickerModel {
	var m ReleasePickerModel
	fillerString := "No Releases Yet..."

	m.common = common
	// TODO: Handle width & height
	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), 40, 20)
	m.list.Title = "Available Releases"
	m.port = viewport.New(40, 20)
	m.port.SetContent(fillerString)
	m.client = picker.GetClient()

	m.examiner = showChangelog // Seems more useful.
	m.changelog = fillerString
	return m
}

func (m ReleasePickerModel) Update(msg tea.Msg) (ReleasePickerModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case NewReleasesMsg:
		// TODO: this logic will need to be reworked if we ever want more releases after init
		m.common.logger.L.Info("recieved releases", "number", len(msg.releases))
		m.ExperimentalReleases = append(m.ExperimentalReleases, msg.releases...)

		// FIXME: is this dance of wrapping and then converting necessary? feels dumb
		relItems := make([]list.Item, len(msg.releases))
		for i := range msg.releases {
			relItems[i] = Release{msg.releases[i]}
		}
		cmd = m.list.SetItems(relItems)

		// Examiner views
		if m.examiner == showRelease {
			firstItem := m.ExperimentalReleases[0]
			m.port.SetContent(releaseView(firstItem))
		}
		var err error
		m.changelog, err = changelogView(msg.releases)
		if err != nil {
			m.common.logger.L.Warn("problem generating changelog", "err", err)
			m.changelog = "Unable to concatenate changelog.\nSee log for details."
		}

	case ErrMsg:
		m.common.logger.L.Error("update got error", "err", msg)
		return m, tea.Quit

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, comKeys.Quit):
			return m, tea.Quit

		case key.Matches(msg, comKeys.Back):
			if !m.list.SettingFilter() {
				// So we don't interfere with list text input
				m.common.Page = Main
			}

		case key.Matches(msg, rKeyMap.toggleExaminer):
			m.examiner = 1 - m.examiner
		}

		// Some updates needed regardless of key pressed
		m.list, cmd = m.list.Update(msg) // Always also pass to list
		// Needs to be updated in case the list moved items
		if m.examiner == showRelease {
			curItem := m.list.SelectedItem().(Release)
			m.port.SetContent(releaseView(curItem.RepositoryRelease))
		} else if m.examiner == showChangelog {
			m.port.SetContent(m.changelog)
		}
	}
	return m, cmd
}

func (m ReleasePickerModel) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Center, m.list.View(), m.port.View())
}

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

func releaseView(r *github.RepositoryRelease) string {
	// TODO: make me pretty
	name := common.ValueOrDefault(r.Name, "Nameless Release")
	tagName := common.ValueOrDefault(r.TagName, "*no tag*")
	body := common.ValueOrDefault(r.Body, "*no body*")
	return fmt.Sprintf("Name: %v\nTag: %v\nBody: %v\n", name, tagName, body)
}

func changelogView(rels []*github.RepositoryRelease) (string, error) {
	// This will break if the body format changes at all, so it must error gracefully
	// TODO: ditto
	var cl strings.Builder
	if rels[len(rels)-1].Name != nil && rels[0].Name != nil {
		fmt.Fprintf(&cl, "Changelog: Releases %v to %v\n", *rels[len(rels)-1].Name, *rels[0].Name)
	} else {
		fmt.Fprintf(&cl, "Changelog:\n")
	}

	for i := range len(rels) {
		if rels[i].Body != nil {
			name := common.ValueOrDefault(rels[i].Name, "Nameless Release")
			cl.WriteString(name)
			cl.WriteString("\n")

			scanner := bufio.NewScanner(strings.NewReader(*rels[i].Body))
			for scanner.Scan() {
				line := scanner.Text()
				if strings.Contains(line, "/pull/") {
					cl.WriteString(line)
				}
			}
			cl.WriteString("\n")
		}
	}
	return cl.String(), nil
}

type releaseKeyMap struct {
	c              commonKeys
	toggleExaminer key.Binding
}

var rKeyMap = releaseKeyMap{
	c:              comKeys,
	toggleExaminer: key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "toggle examine release/changelog")),
}
