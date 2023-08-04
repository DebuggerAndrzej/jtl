package main

import (
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
	"os"
)

var (
	width, _, _        = term.GetSize(int(os.Stdout.Fd()))
	appStyle           = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true)
	itemStyle          = lipgloss.NewStyle().PaddingLeft(4).Width(width - 3)
	selectedItemStyle  = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170")).Width(90)
	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)
