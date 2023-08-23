package ui

import (
	"fmt"
	"strings"
	"time"

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

type finishedProcessing string
type issuesReloadRequired bool

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

func New(jiraClient *jira.Client, config *entities.Config) *Model {
	return &Model{client: jiraClient, config: config}
}

func setIssueListItems(m *Model) error {
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
func (m *Model) enterLoadingScreen() tea.Msg {
	return issuesReloadRequired(true)
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
	err := setIssueListItems(m)
	if err != nil {
		m.actionsHistory += errorLog.Render(
			"Couldn't get issues from Jira API. Check internet connection and vpn if applicable.",
		)
	} else {
		m.actionsHistory += successLog.Render("Initialized JTL")
	}
	m.issueDesc = viewport.New(width/2-7, height-15)
	m.actionsLog = viewport.New(width/2-7, 10)
	setIssueDescription(m)
	m.actionsLog.SetContent(m.actionsHistory)
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

func (m *Model) logHours() tea.Msg {
	var err error
	if timeToLog := m.input.Value(); timeToLog != "" {
		issue_id := getSelectedItemTitle(&m.issues)
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
	issue_id := getSelectedItemTitle(&m.issues)
	status := getSelectedItemStatus(&m.issues)
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
			m.initView(msg.Width, msg.Height-1)
			m.loaded = true
		}
		return m, nil
	case tea.KeyMsg:
		keypress := msg.String()
		if m.input.Focused() {
			switch keypress {
			case "enter":
				if m.input.Value() == "" {
					m.input.Blur()
					return m, nil
				}
				m.loadingText = "Logging hours for issue"

				if m.inputType == "normal" {
					m.actionsHistory += "\n" + infoLog.Render(fmt.Sprintf("Logging  %s hours on %s Issue.", m.input.Value(), getSelectedItemTitle(&m.issues)))
				} else {
					m.actionsHistory += "\n" + infoLog.Render(fmt.Sprintf("Logging  %s hours on %s scrum Issue.", m.input.Value(), getSelectedItemTitle(&m.issues)))
				}

				return m, tea.Sequence(m.enterLoadingScreen, m.logHours)
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
				m.loaded = false
				m.loadingText = "Refreshing list of issues"
				return m, m.dummyRefresh
			case "e":
				m.issueChangeType = "increment"
				m.actionsHistory += "\n" + infoLog.Render(fmt.Sprintf("Incrementing %s status", getSelectedItemTitle(&m.issues)))
				m.loadingText = "Incrementing issues status"
				return m, tea.Sequence(m.enterLoadingScreen, m.changeIssueStatus)
			case "E":
				m.issueChangeType = "decrement"
				m.actionsHistory += "\n" + infoLog.Render(fmt.Sprintf("Decrementing %s status", getSelectedItemTitle(&m.issues)))
				m.loadingText = "Decrementing issues status"
				return m, tea.Sequence(m.enterLoadingScreen, m.changeIssueStatus)
			default:
				m.issues, cmd = m.issues.Update(msg)
				setIssueDescription(&m)
				return m, cmd
			}
		}
	case issuesReloadRequired:
		m.loaded = false
	case finishedProcessing:
		if commandError := string(msg); commandError != "" {
			m.actionsHistory += "\n" + errorLog.Render(commandError)
		} else {
			m.actionsHistory += "\n" + successLog.Render("Last command executed successfully")
			if m.input.Value() != "" {
				logged, _ := time.ParseDuration(m.input.Value())
				m.loggedInSession += logged.Hours()
				loggedInSessionStr := fmt.Sprintf(
					"%sh",
					strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", m.loggedInSession), "0"), "."),
				)
				m.actionsHistory += "\n" + estimateLog.Render(fmt.Sprintf("Logged in session: %s", loggedInSessionStr))
			}
		}
		err := setIssueListItems(&m)
		if err != nil {
			m.actionsHistory += "\n" + errorLog.Render("Couldn't get issues from Jira API. Check internet connection and vpn if applicable.")
		} else {
			m.actionsHistory += "\n" + successLog.Render("Refreshed jira issues")
		}
		m.actionsLog.SetContent(m.actionsHistory)
		m.actionsLog.GotoBottom()
		m.loaded = true
		m.input.Blur()
		m.input.Reset()
	}

	return m, nil
}

func (m Model) View() string {
	if m.loaded {
		baseView := lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.issues.View(),
			lipgloss.JoinVertical(
				lipgloss.Left,
				viewportStyle.Render(m.issueDesc.View()),
				actionsLogStyle.Render(m.actionsLog.View()),
			),
		)
		if m.input.Focused() {
			return appStyle.Render(baseView + "\n " + m.input.View())
		}
		return appStyle.Render(baseView + "\n " + m.help.View(m.keys))
	}
	if m.loadingText == "" {
		return loadingStyle.Render("Loading...")
	}
	return loadingStyle.Render(m.loadingText + "...")

}

func InitTui(config *entities.Config, jiraClient *jira.Client) {
	m := New(jiraClient, config)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
	}
}
