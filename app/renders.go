package app

import (
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	"github.com/dustin/go-humanize"
	"github.com/tubbebubbe/transmission"
)

func RenderTorrentTable(torrents transmission.Torrents, cursor int) []table.Row {
	var tr []table.Row
	for _, t := range torrents {
		strId := strconv.Itoa(t.ID)
		strStatus := TorrentStatus(t)
		strDown := humanize.Bytes(uint64(t.RateDownload))
		strUp := humanize.Bytes(uint64(t.RateUpload))
		tr = append(tr, table.Row{strId, strStatus, strDown, strUp, t.Name})
	}
	return tr
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
