package torrents

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/dustin/go-humanize"
	"github.com/tubbebubbe/transmission"
)

func updateTableRows(m *Model, torrents transmission.Torrents) {
    torrents.SortByAddedDate(true)
	rows := []table.Row{}
	for _, t := range torrents {
		strEta := formatTime(t.Eta)
		strPercentDone := fmt.Sprintf("%.1f", t.PercentDone*100.0)
		strStatus := TorrentStatus(t)
		strDown := humanize.Bytes(uint64(t.RateDownload))
		strUp := humanize.Bytes(uint64(t.RateUpload))
		rows = append(rows, table.Row{strEta, strPercentDone, strStatus, strDown, strUp, t.Name})
	}
	m.table.SetRows(rows)
}

func createTable(rows []table.Row, _ int) table.Model {
	columns := []table.Column{
		{Title: "ETA", Width: 5},
		{Title: "%", Width: 6},
		{Title: "Status", Width: 12},
		{Title: "Down", Width: 6},
		{Title: "Up", Width: 6},
		{Title: "Name", Width: 40},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(6),
	)

	return t
}
