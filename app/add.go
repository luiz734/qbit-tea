package app

import (
	"fmt"
	"log"
	"qbit-tea/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tubbebubbe/transmission"
)

func CmdAdd(m Model) tea.Cmd {
	return func() tea.Msg {
		url := "magnet:?xt=urn:btih:6f1bde857b97b382f8841cdf3a42c530b3f4e34e&dn=archlinux-2024.09.01-x86_64.iso"
		addCommand, err := transmission.NewAddCmdByMagnet(url)
		util.CheckError(err)
		output, err := m.client.ExecuteAddCommand(addCommand)
		log.Println(fmt.Sprintf("%s", output))
		util.CheckError(err)
		return nil
	}
}

type MsgAdd struct{}
