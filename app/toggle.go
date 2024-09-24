package app

import tea "github.com/charmbracelet/bubbletea"

type MsgToggle struct{}

func CmdToggle(m Model) tea.Cmd {
	return func() tea.Msg {
		torrent := m.torrents[m.cursor]
		switch torrent.Status {
		case StatusStopped:
			m.client.StartTorrent(torrent.ID)
		case StatusDownloading, StatusSeeding:
			m.client.StopTorrent(torrent.ID)
		}
		return MsgToggle{}
	}
}
