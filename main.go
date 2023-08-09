package main

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
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
	client *jira.Client
}

func New(jira_client *jira.Client) *Model {
	return &Model{client: jira_client}
}

func setIssueListItems(m *Model) {
	jira_issues := get_all_jira_issues_for_assignee(m.client)
	var s []list.Item
	for _, jira_issue := range jira_issues {
		s = append(s, jira_issue)
	}
	m.issues.ResetSelected()
	m.issues.SetItems(s)
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
	setIssueListItems(m)
}
func getSelectedItemID(l *list.Model) string {
	if i, ok := l.SelectedItem().(Issue); ok {
		return i.title
	} else {
		panic(ok)
	}
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
				time_to_log := m.input.Value()
				issue_id := getSelectedItemID(&m.issues)
				log_hours_for_issue(m.client, issue_id, time_to_log)
				m.input.Blur()
				m.issues.NewStatusMessage(statusMessageStyle(fmt.Sprintf("You logged %s on %s issue", time_to_log, issue_id)))
				setIssueListItems(&m)
			}
			m.input, cmd = m.input.Update(msg)
		}
		if !m.input.Focused() {
			if keypress == "w" {
				m.input.Focus()
			}
			if keypress == "r" {
				setIssueListItems(&m)
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
	return loadingStyle.Render("Loading...")

}

func main() {
	fmt.Println("hej")
	config := get_toml_config()
	client := get_jira_client(config)
	m := New(client)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
	}
}
