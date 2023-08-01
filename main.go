package main

import (
	"fmt"

	list "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	issues list.Model
	err    error
}

func New() *Model {
	return &Model{}
}

func (m *Model) initIssues(width, height int) {
	m.issues = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	m.issues.Title = "Issues"
	m.issues.SetItems([]list.Item{
		Issue{title: "Fake task", short_description: "Some description for this task", status: "Done"},
		Issue{title: "Some task", short_description: "Another description for another task", status: "Done"},
		Issue{title: "Stop messing around", short_description: "Start doing overtimes", status: "Done"},
	})
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.initIssues(msg.Width, msg.Height)
	}
	var cmd tea.Cmd
	m.issues, cmd = m.issues.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.issues.View()
}

func main() {
	m := New()
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
	}
}
