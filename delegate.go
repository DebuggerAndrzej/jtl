package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int  { return 2 }
func (d itemDelegate) Spacing() int { return 1 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	keys := newDelegateKeyMap()
	var title string

	if i, ok := m.SelectedItem().(Issue); ok {
		title = i.Title()
	} else {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.choose):
			return m.NewStatusMessage(statusMessageStyle("You chose " + title))
		case key.Matches(msg, keys.remove):
			index := m.Index()
			m.RemoveItem(index)
			if len(m.Items()) == 0 {
				keys.remove.SetEnabled(false)
			}
			return m.NewStatusMessage(statusMessageStyle("Deleted " + title))
		}
	}
	return nil
}
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Issue)

	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s     %s \n %s", index+1, i.title, i.original_estimate, i.short_description)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type delegateKeyMap struct {
	choose key.Binding
	remove key.Binding
	log    key.Binding
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
		remove: key.NewBinding(
			key.WithKeys("x", "backspace"),
			key.WithHelp("x", "delete"),
		),
		log: key.NewBinding(
			key.WithKeys("w"),
			key.WithHelp("w", "Log hours"),
		),
	}
}
