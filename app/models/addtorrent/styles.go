package addtorrent

import (
	"qbit-tea/app/models"
	"qbit-tea/colors"

	// "github.com/charmbracelet/bubbles/list"
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

var (
	// titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	// paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	// helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	// quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)
