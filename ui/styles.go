package ui

import (
	"os"

	lipgloss "github.com/charmbracelet/lipgloss"
	term "golang.org/x/term"
)

var (
	width, height, _ = term.GetSize(int(os.Stdout.Fd()))
	appStyle         = lipgloss.NewStyle()
	viewportStyle    = lipgloss.NewStyle().
				Width(width/2-7).
				Height(height-15).
				MarginLeft(2).
				Border(lipgloss.RoundedBorder(), true)
	actionsLogStyle = lipgloss.NewStyle().
			Width(width/2-7).
			Height(10).
			MarginLeft(2).
			Border(lipgloss.RoundedBorder(), true)
	itemStyle    = lipgloss.NewStyle().MarginLeft(2).Width(width / 2)
	loadingStyle = lipgloss.NewStyle().
			MarginLeft((width - 10) / 2).
			MarginTop((height - 2) / 2).
			Border(lipgloss.RoundedBorder())
	selectedItemStyle = lipgloss.NewStyle().
				MarginLeft(2).
				Foreground(lipgloss.Color("#03fc6b"))
	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
	successLog        = lipgloss.NewStyle().SetString("SUCCESS: ").Foreground(lipgloss.Color("#AFE1AF"))
	errorLog          = lipgloss.NewStyle().SetString("ERROR: ").Foreground(lipgloss.Color("#FF9999"))
	warningLog        = lipgloss.NewStyle().SetString("WARNING: ").Foreground(lipgloss.Color("#FAD5A5"))
	estimateTimeStyle = lipgloss.NewStyle().Width(20)
	keyStyle          = lipgloss.NewStyle().Width(20)
	statusStyle       = lipgloss.NewStyle().Bold(true)
	rowStyle          = lipgloss.NewStyle().Width(width / 2).PaddingLeft(2)
)
