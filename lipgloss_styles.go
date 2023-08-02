package main

import (
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
	"os"
)

var (
	width, _, _       = term.GetSize(int(os.Stdout.Fd()))
	appStyle          = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).Width(width - 3)
	statusStyle       = lipgloss.NewStyle().Bold(true).PaddingLeft(5)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2).Width(width / 2)
	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(4).
				Foreground(lipgloss.Color("#03fc6b")).
				Width(width / 2)
	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)
