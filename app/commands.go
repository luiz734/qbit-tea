package app

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type MsgToggle struct{}

func CmdToggle(m Model) tea.Cmd {
	return func() tea.Msg {
		torrent := *m.torrent
		switch torrent.Status {
		case StatusStopped:
			m.client.StartTorrent(torrent.ID)
			log.Printf("start torrent %s", torrent.Name)
		case StatusDownloading, StatusSeeding:
			m.client.StopTorrent(torrent.ID)
			log.Printf("stop torrent %s", torrent.Name)
		}
		return MsgToggle{}
	}
}
