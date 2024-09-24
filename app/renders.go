package app

import (
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	"github.com/tubbebubbe/transmission"
)

func RenderTorrentTable(torrents transmission.Torrents, cursor int) []table.Row {
	var tr []table.Row
	for _, t := range torrents {
		strId := strconv.Itoa(t.ID)
		strStatus := TorrentStatus(t)
		tr = append(tr, table.Row{strId, strStatus, t.Name})
	}
	return tr
}
