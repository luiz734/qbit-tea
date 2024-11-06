package torrents

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/tubbebubbe/transmission"
)

func NewAddInDirCmdByMagnet(magnetLink string, path string) (*transmission.Command, error) {
	cmd, _ := transmission.NewAddCmd()
	cmd.Arguments.Filename = magnetLink
	// Can't check if it's a dir on remote hosts
	cmd.Arguments.DownloadDir = path
	return cmd, nil
}

func formatTime(seconds int) string {
	var hours, minutes int
	acc := seconds
	hours = acc / 3600
	acc -= hours * 3600

	minutes = acc / 60
	acc -= minutes * 60

	if hours > 0 {
		return fmt.Sprintf("%dh%dm", hours, minutes)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm", minutes)
	} else {
		return "~"
	}
}

func updateTableSize(t *table.Model, windowSize windowSize) {
	var columns = t.Columns()
	var total int
	for _, c := range columns {
		total += c.Width
	}
	// Last fiels should always be "Name" and fill all the remaining space
	sumBeforeLastColumn := total - columns[len(columns)-1].Width
	nameColumnWidth := windowSize.Width - sumBeforeLastColumn
	columns[len(columns)-1].Width = nameColumnWidth
	t.SetColumns(columns)
	// Let room for other stuff on screen
	t.SetHeight(windowSize.Height - 5)
}
