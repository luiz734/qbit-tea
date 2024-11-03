package errorscreen

import (
	"qbit-tea/app/models"
	"qbit-tea/colors"

	"github.com/charmbracelet/lipgloss"
)

var (
	styleErrTitle = lipgloss.NewStyle().
			Foreground(colors.Pink).
			Bold(true).
			Align(lipgloss.Center)
	styleErrDesc = lipgloss.NewStyle().
			Foreground(colors.Surface2)
	styleError = lipgloss.NewStyle().
		// Margin(1).
		// Border(lipgloss.NormalBorder()).
		// BorderForeground(colors.Surface2).
		Align(lipgloss.Center, lipgloss.Center)

    // Inheritance?
	styleHelp = models.StyleHelp
)
