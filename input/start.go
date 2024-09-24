package input

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

type MsgStart struct {
}

func CmdStart(id int) tea.Cmd {
	return func() tea.Msg {
        idString := strconv.Itoa(id)
		args := []string{"-t", idString, "--start"}
        log.Printf(fmt.Sprintln(id))
		cmd := exec.Command("transmission-remote", args...)
		_, err := cmd.Output()
		if err != nil {
			panic(err)
		}
		// log.Printf(string(output))

		return MsgStart{}
	}
}
