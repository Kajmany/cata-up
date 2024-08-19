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
