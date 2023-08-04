package main

import (
	"fmt"

	list "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	issues list.Model
	err    error
	loaded bool
}

func New() *Model {
	return &Model{}
}

func (m *Model) initIssues(width, height int) {
	m.issues = list.New([]list.Item{}, itemDelegate{}, width, height)
	m.issues.Title = "Issues"
	m.issues.SetItems([]list.Item{
		Issue{title: "Fake task", short_description: "Some description for this task", status: "Done", original_estimate: "2h", logged_time: "0h"},
		Issue{title: "Some task", short_description: "Another description for another task", status: "Done", original_estimate: "4h", logged_time: "2h"},
		Issue{title: "Stop messing around", short_description: "Start doing overtimes", status: "Done", original_estimate: "6h", logged_time: "3h"},
		Issue{title: "Fake task", short_description: "Some description for this task", status: "Done", original_estimate: "2h", logged_time: "0h"},
		Issue{title: "Fake task", short_description: "Some description for this task", status: "Done", original_estimate: "2h", logged_time: "0h"},
		Issue{title: "Fake task", short_description: "Some description for this task", status: "Done", original_estimate: "2h", logged_time: "0h"},
		Issue{title: "Some task", short_description: "Another description for another task", status: "Done", original_estimate: "4h", logged_time: "2h"},
		Issue{title: "Stop messing around", short_description: "Start doing overtimes", status: "Done", original_estimate: "6h", logged_time: "3h"},
		Issue{title: "Some task", short_description: "Another description for another task", status: "Done", original_estimate: "4h", logged_time: "2h"},
		Issue{title: "Stop messing around", short_description: "Start doing overtimes", status: "Done", original_estimate: "6h", logged_time: "3h"},
		Issue{title: "Some task", short_description: "Another description for another task", status: "Done", original_estimate: "4h", logged_time: "2h"},
		Issue{title: "Stop messing around", short_description: "Start doing overtimes", status: "Done", original_estimate: "6h", logged_time: "3h"},
		Issue{title: "Stop messing around", short_description: "Start doing overtimes", status: "Done", original_estimate: "6h", logged_time: "3h"},
		Issue{title: "Stop messing around", short_description: "Start doing overtimes", status: "Done", original_estimate: "6h", logged_time: "3h"},
		Issue{title: "Stop messing around", short_description: "Start doing overtimes", status: "Done", original_estimate: "6h", logged_time: "3h"},
		Issue{title: "Stop messing around", short_description: "Start doing overtimes", status: "Done", original_estimate: "6h", logged_time: "3h"},
		Issue{title: "Some task", short_description: "Another description for another task", status: "Done", original_estimate: "4h", logged_time: "2h"},
		Issue{title: "Stop messing around", short_description: "Start doing overtimes", status: "Done", original_estimate: "6h", logged_time: "3h"},
		Issue{title: "Some task", short_description: "Another description for another task", status: "Done", original_estimate: "4h", logged_time: "2h"},
		Issue{title: "Stop messing around", short_description: "Start doing overtimes", status: "Done", original_estimate: "6h", logged_time: "3h"},
		Issue{title: "Some task", short_description: "Another description for another task", status: "Done", original_estimate: "4h", logged_time: "2h"},
		Issue{title: "Stop messing around", short_description: "Start doing overtimes", status: "Done", original_estimate: "6h", logged_time: "3h"},
		Issue{title: "Some task", short_description: "Another description for another task", status: "Done", original_estimate: "4h", logged_time: "2h"},
		Issue{title: "Stop messing around", short_description: "Start doing overtimes", status: "Done", original_estimate: "6h", logged_time: "3h"},
	})
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			m.initIssues(msg.Width, msg.Height-2)
			m.loaded = true
		}
		return m, nil
	}
	var cmd tea.Cmd
	m.issues, cmd = m.issues.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.loaded {
		return appStyle.Render(m.issues.View())
	}
	return "Loading..."

}

func main() {
	m := New()
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
	}
}
