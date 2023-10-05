package ui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Log             key.Binding
	LogUnderScrum   key.Binding
	IncrementStatus key.Binding
	DecrementStatus key.Binding
	RefreshIssues   key.Binding
	AddIssue        key.Binding
	RemoveIssue     key.Binding
	Quit            key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Log,
		k.LogUnderScrum,
		k.RefreshIssues,
		k.IncrementStatus,
		k.DecrementStatus,
		k.AddIssue,
		k.RemoveIssue,
		k.Quit,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{}, {}}
}

var keys = keyMap{
	Log: key.NewBinding(
		key.WithKeys("w"),
		key.WithHelp("w", "log hours"),
	),
	LogUnderScrum: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "log hours under scrum issue"),
	),
	RefreshIssues: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh issues"),
	),
	IncrementStatus: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "increment issue status"),
	),
	DecrementStatus: key.NewBinding(
		key.WithKeys("E"),
		key.WithHelp("E", "decrement issue status"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	AddIssue: key.NewBinding(
		key.WithKeys("A"),
		key.WithHelp("A", "Add issue by number"),
	),
	RemoveIssue: key.NewBinding(
		key.WithKeys("D"),
		key.WithHelp("D", "Remove selected issue"),
	),
}
