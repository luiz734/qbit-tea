package app

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func viewHeader(m Model) string {
	sectionWidth := m.windowSize.Width / 3
	leftStyle := lipgloss.NewStyle().
		Width(sectionWidth).Align(lipgloss.Left).MarginLeft(1)
	centerStyle := lipgloss.NewStyle().
		Width(sectionWidth).Align(lipgloss.Center).Margin(0)
	rightStyle := lipgloss.NewStyle().
		Width(sectionWidth).Align(lipgloss.Right).MarginRight(1)

	// Left
	var textLeft string
	ratio := 1.0
	if m.torrent != nil {
		ratio = m.torrent.UploadRatio
		strRatio := fmt.Sprintf("%0.1f", ratio)
		textLeft = leftStyle.Render(strRatio)
	}

	// Center
	var styleAddress = lipgloss.NewStyle().
		Align(lipgloss.Center).Bold(true)
	address := styleAddress.Render(fmt.Sprintf("%s", m.address))
	textCenter := centerStyle.Render(address)

	// Right
	r := ":D"
	if ratio < 1.0 {
		r = ":("
	} else if ratio < 2.0 {
		r = ":|"
	}
	textRight := rightStyle.Render(r)

	return lipgloss.JoinHorizontal(lipgloss.Top, textLeft, textCenter, textRight)
}
