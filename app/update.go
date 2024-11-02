package app

import (
	"qbit-tea/util"

	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/log"
	"github.com/tubbebubbe/transmission"
)

func CmdUpdate(m AppModel) tea.Cmd {
	return func() tea.Msg {
		torrents, err := m.client.GetTorrents()
        // log.Info("GET torrent from transmission")
		util.CheckError(err)
		return MsgUpdate{Torrents: torrents}
	}
}

type MsgUpdate struct{ Torrents transmission.Torrents }
