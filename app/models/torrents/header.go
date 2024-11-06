package torrents

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func viewHeader(m Model) string {
	// Left
	var textLeft string
	// Also render when torrent list is empty
	// and there is no ratio to show
	textLeft = leftStyle.Render("--")
	ratio := 1.0
	if m.torrent != nil {
		ratio = m.torrent.UploadRatio
		strRatio := fmt.Sprintf("%0.1f", ratio)
		textLeft = leftStyle.Render(strRatio)
	}

	// Center
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

	headerView := lipgloss.JoinHorizontal(lipgloss.Top, textLeft, textCenter, textRight)
	bottomGapView := strings.Repeat("\n", 2)

	return fmt.Sprint(
		headerView,
		bottomGapView,
	)
}
