package ui

import "github.com/charmbracelet/bubbles/key"

// 'Common' Keys not specific to a single page
// Display methods and full keymap is a per-page struct
// TODO rethink this?

type commonKeys struct {
	Back   key.Binding
	Help   key.Binding
	Quit   key.Binding
	Select key.Binding
}

var comKeys = commonKeys{
	Back: key.NewBinding(
		key.WithKeys("esc", "backspace"),
		key.WithHelp("esc/backspace", "back to Main"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
}

// Front (main) Page
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

// Release Picker Page
type releaseKeyMap struct {
	c              commonKeys
	toggleExaminer key.Binding
	toggleFocus    key.Binding
}

var rKeyMap = releaseKeyMap{
	c:              comKeys,
	toggleExaminer: key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "toggle examine release/changelog")),
	toggleFocus:    key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "toggle list/pager focus")),
}

func (k releaseKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.c.Help, k.c.Quit, k.toggleExaminer, k.toggleFocus}
}

func (k releaseKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.toggleExaminer, k.toggleFocus}, // first column
		{k.c.Help, k.c.Quit},              // second column
	}
}

// Source Picker Page (maybe)
/* Using common until page-specific bindings required
* type sourceKeyMap struct {
*   c commonKeys
*   }
*   */
