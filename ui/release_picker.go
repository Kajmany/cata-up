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

// Couldn't just make it a *tea.Model because they actually aren't.. UGH
type inFocus int

const (
	listFocus inFocus = iota
	portFocus
)

type ReleasePickerModel struct {
	common *Common
	// Will need to swap between fetched release types on request
	ExperimentalReleases []*github.RepositoryRelease
	StableReleases       []list.Item // TODO: not handled
	client               *github.Client
	list                 list.Model
	port                 viewport.Model
	curFocus             inFocus
	examiner             examineMode
	changelog            string

	KeyMap releaseKeyMap
}

func (m ReleasePickerModel) Init() tea.Cmd {
	return m.cmdGetReleases()
}

func NewReleasePicker(common *Common) ReleasePickerModel {
	var m ReleasePickerModel
	fillerString := "No Releases Yet..."

	m.common = common
	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), 20, 40)
	m.list.Title = "Available Releases"
	m.list.SetShowHelp(false)
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
	case tea.WindowSizeMsg:
		m.port.Height = msg.Height
		m.port.Width = msg.Width
		m.list.SetHeight(msg.Height)
		m.list.SetWidth(msg.Width)

	case NewReleasesMsg:
		// TODO: this logic will need to be reworked if we ever want more releases after init
		m.common.logger.L.Info("recieved releases", "number", len(msg.releases))
		m.ExperimentalReleases = append(m.ExperimentalReleases, msg.releases...)

		// TODO: is this dance of wrapping and then converting necessary? feels dumb
		relItems := make([]list.Item, len(msg.releases))
		for i := range msg.releases {
			relItems[i] = Release{msg.releases[i]}
		}
		cmd = m.list.SetItems(relItems)

		// Examiner views
		// Generate changelog regardless
		var err error
		m.changelog, err = changelogView(msg.releases)
		if err != nil {
			m.common.logger.L.Warn("problem generating changelog", "err", err)
			m.changelog = "Unable to concatenate changelog.\nSee log for details."
		}
		// Then pick which to show
		if m.examiner == showRelease {
			firstItem := m.ExperimentalReleases[0]
			m.port.SetContent(releaseView(firstItem))
		} else {
			m.port.SetContent(m.changelog)
		}

	case ErrMsg:
		m.common.logger.L.Error("update got error", "err", msg)
		return m, tea.Quit

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, comKeys.Back):
			if !m.list.SettingFilter() {
				// So we don't interfere with list text input
				m.common.Page = Main
			}

		case key.Matches(msg, rKeyMap.toggleExaminer):
			m.examiner = 1 - m.examiner

		case key.Matches(msg, rKeyMap.toggleFocus):
			m.curFocus = 1 - m.curFocus

		case key.Matches(msg, comKeys.Select):
			// TODO: AHHH!!
		}

		// Regardless of key pressed, curently focused widget also needs it
		if m.curFocus == listFocus {
			m.list, cmd = m.list.Update(msg)
			// Needs to be updated in case the list moved items
			if m.examiner == showRelease {
				curItem := m.list.SelectedItem().(Release)
				m.port.SetContent(releaseView(curItem.RepositoryRelease))
			} else if m.examiner == showChangelog {
				m.port.SetContent(m.changelog)
			}
		} else {
			m.port, cmd = m.port.Update(msg)
		}
	}
	return m, cmd
}

func (m ReleasePickerModel) View() string {
	// TODO: visually distinguish which widget is in focus
	heightBudget := m.common.height
	helpText := m.common.help.View(m)
	heightBudget -= lipgloss.Height(helpText)
	if heightBudget < 10 {
		// Totally arbitrary set point. Mostly concerned about negative.
		m.common.logger.L.Warn("Height available after help text < 10! Expect malformed window.",
			"budget", heightBudget)
	}

	// Help widget width is set in ui, so we need to just worry about content panes here
	widthBudget := (m.common.width / 2)
	// TODO: Do we need a warning about window not wide enough?
	m.list.SetWidth(widthBudget)
	m.port.Width = widthBudget

	// These might be costly calls, but I don't see an alternative
	m.list.SetHeight(heightBudget)
	m.port.Height = heightBudget
	content := lipgloss.JoinHorizontal(lipgloss.Center, m.list.View(), m.port.View())
	return lipgloss.JoinVertical(lipgloss.Center, content, helpText)
}

func (m ReleasePickerModel) ShortHelp() []key.Binding {
	kb := []key.Binding{m.KeyMap.c.Help, m.KeyMap.c.Quit, m.KeyMap.toggleExaminer, m.KeyMap.toggleFocus}
	kb = append(m.list.ShortHelp(), kb...)
	return kb
}

func (m ReleasePickerModel) FullHelp() [][]key.Binding {
	kbs := [][]key.Binding{
		{m.KeyMap.toggleExaminer, m.KeyMap.toggleFocus}, // first column
		{m.KeyMap.c.Help, m.KeyMap.c.Quit},              // second column
	}
	kbs = append(m.list.FullHelp(), kbs...)
	return kbs
}

func releaseView(r *github.RepositoryRelease) string {
	// TODO: make me pretty
	name := common.ValueOrDefault(r.Name, "Nameless Release")
	tagName := common.ValueOrDefault(r.TagName, "*no tag*")
	body := common.ValueOrDefault(r.Body, "*no body*")
	return fmt.Sprintf("Name: %v\nTag: %v\nBody: %v\n", name, tagName, body)
}

// The special changelog algorithm sauce? Check every line. Both C:DDA and C:BN
// put github PR links on each change item line. Definitely won't work everywhere.
func changelogView(rels []*github.RepositoryRelease) (string, error) {
	// TODO: ditto
	var cl strings.Builder
	if rels[len(rels)-1].Name != nil && rels[0].Name != nil {
		fmt.Fprintf(&cl, "Changelog: Releases %v to %v\n", *rels[len(rels)-1].Name, *rels[0].Name)
	} else {
		fmt.Fprintf(&cl, "Changelog:\n")
	}

	for i := range len(rels) {
		// Did we have a line-item written (or is filler-text needed?)
		wrote := false

		name := common.ValueOrDefault(rels[i].Name, "Nameless Release")
		cl.WriteString(name)
		cl.WriteString("\n")
		if rels[i].Body != nil {
			scanner := bufio.NewScanner(strings.NewReader(*rels[i].Body))
			for scanner.Scan() {
				line := scanner.Text()
				if strings.Contains(line, "/pull/") {
					cl.WriteString(line)
					wrote = true
				}
			}
			cl.WriteString("\n")
		}
		if !wrote {
			cl.WriteString("**no line-item changes extracted**\n")
		}
	}
	return cl.String(), nil
}
