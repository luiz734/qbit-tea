package app

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tubbebubbe/transmission"
	"log"
	"qbit-tea/util"
)

func CmdRemove(m Model, deleteData bool) tea.Cmd {
	return func() tea.Msg {
		torrent := m.torrents[m.cursor]
		deleteCommand, err := transmission.NewDelCmd(torrent.ID, deleteData)
		output, err := m.client.ExecuteCommand(deleteCommand)
		log.Println(fmt.Sprintf("%s", output))
		util.CheckError(err)

		return nil
	}
}
