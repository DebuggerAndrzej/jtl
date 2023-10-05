package ui

import (
	"fmt"
	"strings"
	"time"

	jira "github.com/andygrunwald/go-jira"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"

	"github.com/DebuggerAndrzej/jtl/backend/entities"
)

type successProcessing string
type failedProcessing string
type issuesReloadRequired bool

func New(jiraClient *jira.Client, config *entities.Config, configPath string) *Model {
	return &Model{client: jiraClient, config: config, configPath: configPath}
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
				if m.inputType == "add issue" {
					m.loadingText = fmt.Sprintf("Adding issue with id %s", m.input.Value())
					return m, tea.Sequence(m.enterLoadingScreen, m.addIssue)
				}

				m.loadingText = "Logging hours for issue"
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
				m.loadingText = "Incrementing issues status"
				return m, tea.Sequence(m.enterLoadingScreen, m.changeIssueStatus)
			case "E":
				m.issueChangeType = "decrement"
				m.loadingText = "Decrementing issues status"
				return m, tea.Sequence(m.enterLoadingScreen, m.changeIssueStatus)
			case "A":
				m.inputType = "add issue"
				m.input.Focus()
			case "D":
				return m, tea.Sequence(m.enterLoadingScreen, m.removeIssue)
			default:
				m.issues, cmd = m.issues.Update(msg)
				m.setIssueDescription()
				return m, cmd
			}
		}

	case issuesReloadRequired:
		m.loaded = false

	case successProcessing:
		if sucessMsg := string(msg); sucessMsg != "" {
			m.actionsHistory += "\n" + successLog.Render(sucessMsg)
			if m.input.Value() != "" {
				logged, _ := time.ParseDuration(m.input.Value())
				m.loggedInSession += logged.Hours()
				loggedInSessionStr := fmt.Sprintf(
					"%sh",
					strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", m.loggedInSession), "0"), "."),
				)
				if loggedInSessionStr == "8h" {
					m.hoursSummary.SetContent(loggedStyle.Render(loggedTimeStyle.Render("Great job, see you tomorrow! :)")))
				} else {
					m.hoursSummary.SetContent(loggedStyle.Render(fmt.Sprintf("Logged in session: %s", loggedTimeStyle.Render(loggedInSessionStr))))
				}
			}
		}
		if err := m.setIssueListItems(); err != nil {
			m.actionsHistory += "\n" + errorLog.Render("Couldn't get issues from Jira API. Check internet connection and vpn if applicable.")
		} else {
			m.actionsHistory += "\n" + successLog.Render("Refreshed jira issues")
		}
		m.actionsLog.SetContent(m.actionsHistory)
		m.actionsLog.GotoBottom()
		m.loaded = true
		m.input.Blur()
		m.input.Reset()

	case failedProcessing:
		m.actionsHistory += "\n" + errorLog.Render(string(msg))
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
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.hoursSummary.View(),
				m.issues.View(),
			),
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

func InitTui(config *entities.Config, jiraClient *jira.Client, configPath string) {
	m := New(jiraClient, config, configPath)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
	}
}
