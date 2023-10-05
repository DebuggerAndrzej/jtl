package ui

import (
	"fmt"
	"strconv"

	jira "github.com/andygrunwald/go-jira"
	help "github.com/charmbracelet/bubbles/help"
	list "github.com/charmbracelet/bubbles/list"
	textinput "github.com/charmbracelet/bubbles/textinput"
	viewport "github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	glamour "github.com/charmbracelet/glamour"

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
	configPath      string
}

func (m *Model) enterLoadingScreen() tea.Msg {
	return issuesReloadRequired(true)
}

func (m *Model) addIssue() tea.Msg {
	_, err := strconv.ParseInt(m.input.Value(), 10, 64)

	if err != nil {
		return failedProcessing("Please provide a number.")
	}
	backend.AddIssueToConfig(m.configPath, m.input.Value(), m.config)
	return successProcessing(fmt.Sprintf("Issue %s added to config", m.input.Value()))

}
func (m *Model) removeIssue() tea.Msg {
	issueId := m.getSelectedItemTitle()
	backend.RemoveIssueFromConfig(m.configPath, issueId, m.config)
	return successProcessing(fmt.Sprintf("Issue %s removed from config", issueId))

}

func (m *Model) logHours() tea.Msg {
	var err error
	var successMsg string
	if timeToLog := m.input.Value(); timeToLog != "" {
		issueId := m.getSelectedItemTitle()
		if m.inputType == "normal" {
			successMsg = fmt.Sprintf("Logged  %s hours on %s Issue.", m.input.Value(), m.getSelectedItemTitle())
			err = backend.LogHoursForIssue(m.client, issueId, timeToLog)
		} else {
			successMsg = fmt.Sprintf("Logged  %s hours on %s Scrum Issue.", m.input.Value(), m.getSelectedItemTitle())
			err = backend.LogHoursForIssuesScrumMeetings(m.client, issueId, timeToLog)
		}
	}

	if err != nil {
		return failedProcessing(err.Error())
	}

	return successProcessing(successMsg)
}

func (m *Model) changeIssueStatus() tea.Msg {
	var err error
	var successMsg string
	issue_id := m.getSelectedItemTitle()
	status := m.getSelectedItemStatus()
	if m.issueChangeType == "increment" {
		successMsg = fmt.Sprintf("Incremented %s status", m.getSelectedItemTitle())
		err = backend.IncrementIssueStatus(m.client, issue_id, status)
	} else {
		successMsg = fmt.Sprintf("Decremented %s status", m.getSelectedItemTitle())
		err = backend.DecrementIssueStatus(m.client, issue_id, status)
	}

	if err != nil {
		return failedProcessing(err.Error())
	}
	return successProcessing(successMsg)
}

func (m *Model) dummyRefresh() tea.Msg {
	return successProcessing("")
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
	m.input = input
	m.issues = list.New([]list.Item{}, itemDelegate{}, width/2, height-5)
	m.issues.SetShowHelp(false)
	m.issues.SetShowStatusBar(false)
	m.issues.SetFilteringEnabled(false)
	m.issues.SetShowTitle(false)
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
	m.hoursSummary = viewport.New((width/2)+2, 4)
	m.hoursSummary.SetContent(loggedStyle.Render(fmt.Sprintf("Logged in session: %s", loggedTimeStyle.Render("0h"))))
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
