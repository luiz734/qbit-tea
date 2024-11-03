package addtorrent

import (
	"qbit-tea/app/models"
	"qbit-tea/colors"

	"github.com/charmbracelet/lipgloss"
)

var (
	styleLabel = lipgloss.NewStyle().
			Bold(true).
			Foreground(colors.Pink)

	styleMagnet = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colors.Surface0)

	styleHelp = models.StyleHelp
)
