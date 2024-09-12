package ui

import (
	"github.com/charmbracelet/bubbles/key"
)

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
		key.WithHelp("s", "source selection"),
	),
	GetReleases: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "get releases"),
	),
	Channel: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "change release channel"),
	),
}

// Release Picker Page
type releaseKeyMap struct {
	c              commonKeys
	toggleExaminer key.Binding
	toggleFocus    key.Binding
}

// TODO: It's not shown here that this is changed in-flight to include list keymaps on display
var rKeyMap = releaseKeyMap{
	c:              comKeys,
	toggleExaminer: key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "toggle examine release/changelog")),
	toggleFocus:    key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "toggle list/pager focus")),
}

// TODO: Finish moving help funcs/ INITs to respective pages (more dynamic logic needs the parent model!!)

// Source Picker Page
type sourceKeyMap struct {
	c commonKeys
}
