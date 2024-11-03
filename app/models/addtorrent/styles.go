package addtorrent

import (
	"qbit-tea/app/models"
	"qbit-tea/colors"

	// "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	stylePrompt = lipgloss.NewStyle().
			Bold(true).
			Foreground(colors.Pink)
	styleUnfocused = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colors.Surface0).
			MarginTop(1)
		//          PaddingRight(1)

	styleFocused = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colors.Pink).
			MarginTop(1)
		// PaddingRight(1)

	styleHelp = models.StyleHelp
)

// Directory picker
var (
	titleStyle = lipgloss.NewStyle().
			MarginLeft(0).
			Foreground(colors.Pink).
			Bold(true)
	itemStyle = lipgloss.NewStyle().
			PaddingLeft(4).
			Foreground(colors.Surface2)
	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("#fff"))

	// paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	// helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	// quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)
