package ui

import (
	jira "github.com/andygrunwald/go-jira"
	help "github.com/charmbracelet/bubbles/help"
	list "github.com/charmbracelet/bubbles/list"
	textinput "github.com/charmbracelet/bubbles/textinput"
	viewport "github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	glamour "github.com/charmbracelet/glamour"
	lipgloss "github.com/charmbracelet/lipgloss"

	"github.com/DebuggerAndrzej/jtl/backend"
	"github.com/DebuggerAndrzej/jtl/backend/entities"
)

type Model struct {
	issues          list.Model
	input           textinput.Model
	issueDesc       viewport.Model
	actionsLog      viewport.Model
	hoursSummary    viewport.Model
	err             error
	loaded          bool
	help            help.Model
	keys            keyMap
	client          *jira.Client
	config          *entities.Config
	inputType       string
	loadingText     string
	issueChangeType string
	actionsHistory  string
	loggedInSession float64
}

func (m *Model) enterLoadingScreen() tea.Msg {
	return issuesReloadRequired(true)
}

func (m *Model) logHours() tea.Msg {
	var err error
	if timeToLog := m.input.Value(); timeToLog != "" {
		issue_id := m.getSelectedItemTitle()
		if m.inputType == "normal" {
			err = backend.LogHoursForIssue(m.client, issue_id, timeToLog)
		} else {
			err = backend.LogHoursForIssuesScrumMeetings(m.client, issue_id, timeToLog)
		}
	}

	if err != nil {
		return finishedProcessing(err.Error())
	}
	return finishedProcessing("")
}

func (m *Model) changeIssueStatus() tea.Msg {
	var err error
	issue_id := m.getSelectedItemTitle()
	status := m.getSelectedItemStatus()
	if m.issueChangeType == "increment" {
		err = backend.IncrementIssueStatus(m.client, issue_id, status)
	} else {
		err = backend.DecrementIssueStatus(m.client, issue_id, status)
	}

	if err != nil {
		return finishedProcessing(err.Error())
	}
	return finishedProcessing("")
}

func (m *Model) dummyRefresh() tea.Msg {
	return finishedProcessing("")
}

func (m *Model) setIssueListItems() error {
	jiraIssues, err := backend.GetAllJiraIssuesForAssignee(m.client, m.config)
	if err != nil {
		return err
	}
	var issues []list.Item
	for _, jiraIssue := range jiraIssues {
		issues = append(issues, jiraIssue)
	}
	m.issues.SetItems(issues)
	return nil
}

func (m *Model) initView(width, height int) {
	m.keys = keys
	m.help = help.New()
	input := textinput.New()
	input.Placeholder = "Log hours in (float)h format"
	m.input = input
	m.issues = list.New([]list.Item{}, itemDelegate{}, width, height)
	m.issues.SetShowHelp(false)
	m.issues.SetShowStatusBar(false)
	m.issues.SetFilteringEnabled(false)
	m.issues.Styles.Title = lipgloss.NewStyle()
	m.issues.Title = ""
	err := m.setIssueListItems()
	if err != nil {
		m.actionsHistory += errorLog.Render(
			"Couldn't get issues from Jira API. Check internet connection and vpn if applicable.",
		)
	} else {
		m.actionsHistory += successLog.Render("Initialized JTL")
	}
	m.issueDesc = viewport.New(width/2-7, height-15)
	m.actionsLog = viewport.New(width/2-7, 10)
	m.setIssueDescription()
	m.actionsLog.SetContent(m.actionsHistory)
}
func (m *Model) getSelectedItemTitle() string {
	if i, ok := m.issues.SelectedItem().(entities.Issue); ok {
		return i.Key
	} else {
		panic(ok)
	}
}

func (m *Model) getSelectedItemStatus() string {
	if i, ok := m.issues.SelectedItem().(entities.Issue); ok {
		return i.Status
	} else {
		panic(ok)
	}
}

func (m *Model) setIssueDescription() {
	var desc string
	if i, ok := m.issues.SelectedItem().(entities.Issue); ok {
		desc = i.Description
	}
	out, _ := glamour.Render(desc, "dark")
	m.issueDesc.SetContent(out)
}
