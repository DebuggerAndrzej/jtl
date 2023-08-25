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

type finishedProcessing string
type issuesReloadRequired bool

func New(jiraClient *jira.Client, config *entities.Config) *Model {
	return &Model{client: jiraClient, config: config}
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
					m.actionsHistory += "\n" + infoLog.Render(fmt.Sprintf("Logging  %s hours on %s Issue.", m.input.Value(), m.getSelectedItemTitle()))
				} else {
					m.actionsHistory += "\n" + infoLog.Render(fmt.Sprintf("Logging  %s hours on %s scrum Issue.", m.input.Value(), m.getSelectedItemTitle()))
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
				m.actionsHistory += "\n" + infoLog.Render(fmt.Sprintf("Incrementing %s status", m.getSelectedItemTitle()))
				m.loadingText = "Incrementing issues status"
				return m, tea.Sequence(m.enterLoadingScreen, m.changeIssueStatus)
			case "E":
				m.issueChangeType = "decrement"
				m.actionsHistory += "\n" + infoLog.Render(fmt.Sprintf("Decrementing %s status", m.getSelectedItemTitle()))
				m.loadingText = "Decrementing issues status"
				return m, tea.Sequence(m.enterLoadingScreen, m.changeIssueStatus)
			default:
				m.issues, cmd = m.issues.Update(msg)
				m.setIssueDescription()
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
		err := m.setIssueListItems()
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
