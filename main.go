package main

import (
	"fmt"

	help "github.com/charmbracelet/bubbles/help"
	list "github.com/charmbracelet/bubbles/list"
	textinput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	issues list.Model
	input  textinput.Model
	err    error
	loaded bool
	help   help.Model
	keys   keyMap
}

func New() *Model {
	return &Model{}
}

func (m *Model) initIssues(width, height int) {
	m.keys = keys
	m.help = help.New()
	input := textinput.New()
	input.Placeholder = "Log hours in (float)h format"
	input.CharLimit = 250
	input.Width = 50
	m.input = input
	m.issues = list.New([]list.Item{}, itemDelegate{}, width, height)
	m.issues.Title = "Issues"
	m.issues.SetShowHelp(false)
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
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			m.initIssues(msg.Width, msg.Height-3)
			m.loaded = true
		}
		return m, nil
	case tea.KeyMsg:
		keypress := msg.String()
		if m.input.Focused() {
			if keypress == "q" {
				return m, tea.Quit
			}
			if keypress == "enter" {
				_ = m.input.Value()
			}
			m.input, cmd = m.input.Update(msg)
		}
		if !m.input.Focused() {
			if keypress == "w" {
				m.input.Focus()
			}
		}

	}
	if !m.input.Focused() {
		m.issues, cmd = m.issues.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) View() string {
	if m.loaded {
		if m.input.Focused() {
			return appStyle.Render(m.issues.View() + "\n" + m.input.View())
		}
		return appStyle.Render(m.issues.View() + "\n" + "  " + m.help.View(m.keys))
	}
	return "Loading..."

}

func main() {
	_ = get_toml_config()
	m := New()
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
	}
}
