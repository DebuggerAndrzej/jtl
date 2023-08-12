package ui

import (
	"fmt"

	jira "github.com/andygrunwald/go-jira"
	help "github.com/charmbracelet/bubbles/help"
	list "github.com/charmbracelet/bubbles/list"
	textinput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/amazurki/JTL/backend"
	"github.com/amazurki/JTL/backend/entities"
)

type Model struct {
	issues     list.Model
	input      textinput.Model
	err        error
	loaded     bool
	help       help.Model
	keys       keyMap
	client     *jira.Client
	config     *Config
	input_type string
}

func New(jiraClient *jira.Client, config *Config) *Model {
	return &Model{client: jiraClient, config: config}
}

func setIssueListItems(m *Model) {
	jiraIssues := getAllJiraIssuesForAssignee(m.client, m.config)
	var issues []list.Item
	for _, jiraIssue := range jiraIssues {
		issues = append(issues, jiraIssue)
	}
	m.issues.ResetSelected()
	m.issues.SetItems(issues)
}

func (m *Model) initView(width, height int) {
	m.keys = keys
	m.help = help.New()
	input := textinput.New()
	input.Placeholder = "Log hours in (float)h format"
	m.input = input
	m.issues = list.New([]list.Item{}, itemDelegate{}, width, height)
	m.issues.Title = "Issues"
	m.issues.SetShowHelp(false)
	setIssueListItems(m)
}

func getSelectedItemTitle(l *list.Model) string {
	if i, ok := l.SelectedItem().(Issue); ok {
		return i.title
	} else {
		panic(ok)
	}
}

func getSelectedItemStatus(l *list.Model) string {
	if i, ok := l.SelectedItem().(Issue); ok {
		return i.status
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
			m.initView(msg.Width, msg.Height-3)
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
				if time_to_log != "" {
					issue_id := getSelectedItemTitle(&m.issues)
					if m.input_type == "normal" {
						log_hours_for_issue(m.client, issue_id, time_to_log)
						m.issues.NewStatusMessage(statusMessageStyle(fmt.Sprintf("You logged %s on %s issue", time_to_log, issue_id)))
						setIssueListItems(&m)
					} else {
						logHoursForIssuesScrumMeetings(m.client, issue_id, time_to_log)
						m.issues.NewStatusMessage(statusMessageStyle(fmt.Sprintf("You logged %s on %s issue's scrum meetings", time_to_log, issue_id)))
					}

				}
				m.input.Blur()
			}
			m.input, cmd = m.input.Update(msg)
		}
		if !m.input.Focused() {
			if keypress == "w" {
				m.input_type = "normal"
				m.input.Focus()
			}
			if keypress == "s" {
				m.input_type = "scrum"
				m.input.Focus()
			}
			if keypress == "r" {
				setIssueListItems(&m)
			}
			if keypress == "e" {
				issue_id := getSelectedItemTitle(&m.issues)
				status := getSelectedItemStatus(&m.issues)
				incrementIssueStatus(m.client, issue_id, status)
				m.issues.NewStatusMessage(statusMessageStyle(fmt.Sprintf("You incremented status on %s issue", issue_id)))
				setIssueListItems(&m)
			}
			if keypress == "E" {
				issue_id := getSelectedItemTitle(&m.issues)
				status := getSelectedItemStatus(&m.issues)
				decrementIssueStatus(m.client, issue_id, status)
				m.issues.NewStatusMessage(statusMessageStyle(fmt.Sprintf("You decremented status on %s issue", issue_id)))
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

func initTui(jiraClient *jira.Client, config *Config) {
	m := New(jiraClient, config)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
	}
}
