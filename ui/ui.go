package ui

import (
	"fmt"

	jira "github.com/andygrunwald/go-jira"
	help "github.com/charmbracelet/bubbles/help"
	list "github.com/charmbracelet/bubbles/list"
	textinput "github.com/charmbracelet/bubbles/textinput"
	viewport "github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	glamour "github.com/charmbracelet/glamour"
	lipgloss "github.com/charmbracelet/lipgloss"

	"jtl/backend"
	"jtl/backend/entities"
)

type Model struct {
	issues    list.Model
	input     textinput.Model
	issueDesc viewport.Model
	err       error
	loaded    bool
	help      help.Model
	keys      keyMap
	client    *jira.Client
	config    *entities.Config
	inputType string
}

func New(jiraClient *jira.Client, config *entities.Config) *Model {
	return &Model{client: jiraClient, config: config}
}

func setIssueListItems(m *Model) {
	jiraIssues := backend.GetAllJiraIssuesForAssignee(m.client, m.config)
	var issues []list.Item
	for _, jiraIssue := range jiraIssues {
		issues = append(issues, jiraIssue)
	}
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
	m.issueDesc = viewport.New(width/2-10, height-21)
	setIssueDescription(m)
}

func getSelectedItemTitle(l *list.Model) string {
	if i, ok := l.SelectedItem().(entities.Issue); ok {
		return i.Key
	} else {
		panic(ok)
	}
}

func getSelectedItemStatus(l *list.Model) string {
	if i, ok := l.SelectedItem().(entities.Issue); ok {
		return i.Status
	} else {
		panic(ok)
	}
}

func logHours(m *Model) tea.Cmd {
	if timeToLog := m.input.Value(); timeToLog != "" {
		issue_id := getSelectedItemTitle(&m.issues)
		if m.inputType == "normal" {
			backend.LogHoursForIssue(m.client, issue_id, timeToLog)
			m.issues.NewStatusMessage(statusMessageStyle(fmt.Sprintf("You logged %s on %s issue", timeToLog, issue_id)))
			setIssueListItems(m)
		} else {
			backend.LogHoursForIssuesScrumMeetings(m.client, issue_id, timeToLog)
			m.issues.NewStatusMessage(statusMessageStyle(fmt.Sprintf("You logged %s on %s issue's scrum meetings", timeToLog, issue_id)))
		}
	}
	m.input.Blur()
	m.input.Reset()
	_, cmd := m.input.Update(nil)
	return cmd
}

func changeIssueStatus(m *Model, changeType string) {
	issue_id := getSelectedItemTitle(&m.issues)
	status := getSelectedItemStatus(&m.issues)
	if changeType == "increment" {
		backend.IncrementIssueStatus(m.client, issue_id, status)
	} else {
		backend.DecrementIssueStatus(m.client, issue_id, status)
	}
	m.issues.NewStatusMessage(statusMessageStyle(fmt.Sprintf("You %sed status on %s issue", changeType, issue_id)))
	setIssueListItems(m)
}

func setIssueDescription(m *Model) {
	var desc string
	if i, ok := m.issues.SelectedItem().(entities.Issue); ok {
		desc = i.Description
	}
	out, _ := glamour.Render(desc, "dark")
	m.issueDesc.SetContent(out)
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
			switch keypress {
			case "enter":
				cmd = logHours(&m)
				return m, cmd
			}
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		} else {
			switch keypress {
			case "w":
				m.inputType = "normal"
				m.input.Focus()
			case "s":
				m.inputType = "scrum"
				m.input.Focus()
			case "r":
				setIssueListItems(&m)
			case "e":
				changeIssueStatus(&m, "increment")
			case "E":
				changeIssueStatus(&m, "decrement")
			default:
				m.issues, cmd = m.issues.Update(msg)
				setIssueDescription(&m)
				return m, cmd
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.loaded {
		if m.input.Focused() {
			return appStyle.Render(
				lipgloss.JoinHorizontal(
					lipgloss.Left,
					m.issues.View(),
					viewportStyle.Render(m.issueDesc.View()),
				) + "\n" + m.input.View(),
			)
		}
		return appStyle.Render(
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				m.issues.View(),
				viewportStyle.Render(m.issueDesc.View()),
			) + "\n" + "  " + m.help.View(
				m.keys,
			),
		)
	}
	return loadingStyle.Render("Loading...")

}

func InitTui(config *entities.Config, jiraClient *jira.Client) {
	m := New(jiraClient, config)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
	}
}
