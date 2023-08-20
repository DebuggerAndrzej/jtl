package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/DebuggerAndrzej/jtl/backend/entities"
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 2 }
func (d itemDelegate) Spacing() int                              { return 1 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(entities.Issue)

	if !ok {
		return
	}
	estimateTime := estimateTimeStyle.Render(fmt.Sprintf("%s of %s", i.LoggedTime, i.OriginalEstimate))
	firstRow := rowStyle.Render(keyStyle.Render(i.Key) + estimateTime + statusStyle.Render(i.Status))
	secondRow := rowStyle.Render(i.ShortDescription)

	singleIssue := lipgloss.JoinVertical(lipgloss.Left, firstRow, secondRow)

	if index == m.Index() {
		fmt.Fprint(w, selectedItemStyle.Render(singleIssue))
	} else {
		fmt.Fprint(w, singleIssue)
	}

}
