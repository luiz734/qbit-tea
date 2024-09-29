package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
)

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

func UpdateColumnsWidth(t table.Model, windowWidth int) table.Model {
	var columns = t.Columns()
	var total int
	for _, c := range columns {
		total += c.Width
	}
	// Last fiels should always be "Name" and fill all the remaining space
	sumBeforeLastColumn := total - columns[len(columns)-1].Width
	nameColumnWidth := windowWidth - sumBeforeLastColumn
	columns[len(columns)-1].Width = nameColumnWidth
	t.SetColumns(columns)
	return t
}
