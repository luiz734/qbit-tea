package app

import (
	"fmt"
	"log"
	"qbit-tea/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tubbebubbe/transmission"
)

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
		return nil
	}
}

func CmdRemove(m Model, deleteData bool) tea.Cmd {
	return func() tea.Msg {
		torrent := *m.torrent
		deleteCommand, err := transmission.NewDelCmd(torrent.ID, deleteData)
		log.Printf("delete torrent %s", torrent.Name)
		output, err := m.client.ExecuteCommand(deleteCommand)
		log.Println(fmt.Sprintf("%s", output))
		util.CheckError(err)
		return nil
	}
}
