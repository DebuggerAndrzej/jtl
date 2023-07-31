package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	cursor int
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	}
	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("Hi. This program won't exit. To quit press any key.\n")
}

func main() {
	p := tea.NewProgram(model{5}, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
	}
}
