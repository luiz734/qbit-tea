package torrents

import (
	"qbit-tea/app/models"

	"github.com/charmbracelet/lipgloss"
)

// Header
var (
	leftStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			MarginLeft(1)
	centerStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Margin(0)
	rightStyle = lipgloss.NewStyle().
			Align(lipgloss.Right).
			MarginRight(1)

	styleAddress = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Bold(true)
)

var styleHelp = models.StyleHelp
var styleEmptyTorrent = styleHelp
