package app

import (
	"qbit-tea/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tubbebubbe/transmission"
)

func CmdUpdate(m Model) tea.Cmd {
	return func() tea.Msg {
		torrents, err := m.client.GetTorrents()
        // log.Print(fmt.Sprintf("%s", torrents))
		util.CheckError(err)
		return MsgUpdate{Torrents: torrents}
	}
}

type MsgUpdate struct{ Torrents transmission.Torrents }
