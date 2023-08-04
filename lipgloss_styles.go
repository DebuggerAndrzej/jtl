package main

import "github.com/charmbracelet/lipgloss"

var (
	appStyle           = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).Padding(1, 5, 1, 5)
	itemStyle          = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle  = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)
