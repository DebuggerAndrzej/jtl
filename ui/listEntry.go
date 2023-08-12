package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"jtl/backend/entities"
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

	str := fmt.Sprintf(
		"%d. %s    logged %s of %s estimate  %s \n %s",
		index+1,
		i.Key,
		i.LoggedTime,
		i.OriginalEstimate,
		statusStyle.Render(i.Status),
		i.ShortDescription,
	)

	if index == m.Index() {
		fmt.Fprint(w, selectedItemStyle.Render(str))
	} else {
		fmt.Fprint(w, itemStyle.Render(str))
	}

}
