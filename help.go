package main

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Log           key.Binding
	RefreshIssues key.Binding
	Filter        key.Binding
	Quit          key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Log, k.RefreshIssues, k.Filter, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{}, {}}
}

var keys = keyMap{
	Log: key.NewBinding(
		key.WithKeys("w"),
		key.WithHelp("w", "log hours"),
	),
	RefreshIssues: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh issues"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter issues"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
