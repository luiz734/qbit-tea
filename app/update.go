package app

import (
	"log"
	"qbit-tea/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tubbebubbe/transmission"
)

func CmdUpdate(m Model) tea.Cmd {
	return func() tea.Msg {
		torrents, err := m.client.GetTorrents()
        log.Print("GET torrent from transmission")
		util.CheckError(err)
		return MsgUpdate{Torrents: torrents}
	}
}

type MsgUpdate struct{ Torrents transmission.Torrents }
