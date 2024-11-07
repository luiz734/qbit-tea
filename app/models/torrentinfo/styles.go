package torrentinfo

import (
	"qbit-tea/app/models"
	"qbit-tea/colors"

	"github.com/charmbracelet/lipgloss"
)

var (
	styleHelp  = models.StyleHelp
	styleLabel = lipgloss.NewStyle().
			Bold(true).
			Foreground(colors.Pink)
	styleValue = lipgloss.NewStyle().
			Foreground(colors.Surface2)
	styleSeparator = lipgloss.NewStyle().
			Foreground(colors.Surface0)

	styleWrapper = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(colors.Surface0).
            Padding(1, 3)
)
