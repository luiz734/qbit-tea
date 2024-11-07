package torrents

import (
	"fmt"
	"qbit-tea/app/models"

	"github.com/charmbracelet/log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tubbebubbe/transmission"
)

func CmdToggle(m Model) tea.Cmd {
	return func() tea.Msg {
		torrent := *m.torrent
		switch torrent.Status {
		case models.StatusStopped:
			m.client.StartTorrent(torrent.ID)
			log.Printf("start torrent %s", torrent.Name)
		case models.StatusDownloading, models.StatusSeeding:
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
		log.Printf(fmt.Sprintf("%v", output))
		if err != nil {
			panic(err)
		}
		// util.CheckError(err)
		return nil
	}
}

type MsgUpdate struct{ Torrents transmission.Torrents }
type MsgError struct {
	title     string
	err       error
	prevModel tea.Model
}

func CmdUpdate(m Model) tea.Cmd {
	torrents, err := m.client.GetTorrents()
	if err != nil {
		log.Error("Can't update: %v", err)
		return func() tea.Msg {
			return MsgError{
				title:     "Can't reach transmission-daemon",
				err:       err,
				prevModel: models.GetQuitModel(),
			}
		}
	}
	return func() tea.Msg {
		return MsgUpdate{Torrents: torrents}
	}

}
