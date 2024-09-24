package input

import (
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type UserAction struct {
	Shortcuts   []string
	Description string
	Command     tea.Cmd
}

var actionsMap = []UserAction{
	{[]string{"a"}, "Add", CmdAdd},
	// {[]string{"u"}, "Update", CmdUpdate},
	{[]string{"k", "up"}, "Move Up", CmdMoveCursorUp},
	{[]string{"j", "down"}, "Move Down", CmdMoveCursorDown},
	// {[]string{"p"}, "Pause/Unpause", CmdStart()},
	{[]string{"d"}, "Delete", CmdDefault},
	{[]string{"x"}, "Delete and remove files", CmdDefault},
	{[]string{"q"}, "Quit", tea.Quit},
}

func HelpMsg() string {

	var output strings.Builder
	output.WriteString("---\n")
	for _, action := range actionsMap {
		shortcuts := strings.Join(action.Shortcuts, ", ")
		output.WriteString(fmt.Sprintf("%-8s -> %s\n", shortcuts, action.Description))
	}
	return output.String()
}

type MsgDefault struct{}

type MsgMoveCursor struct {
	Movement int
}

type MsgAdd struct {
	Foo string
}

func CmdMoveCursorUp() tea.Msg {
	return MsgMoveCursor{-1}
}
func CmdMoveCursorDown() tea.Msg {
	return MsgMoveCursor{1}
}

func CmdAdd() tea.Msg {
	// transmission-remote --add "url"
	url := "magnet:?xt=urn:btih:6f1bde857b97b382f8841cdf3a42c530b3f4e34e&dn=archlinux-2024.09.01-x86_64.iso"
	args := []string{"-a", url}
	cmd := exec.Command("transmission-remote", args...)
	stdout, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return MsgAdd{string(stdout)}
}

func CmdDefault() tea.Msg {
	return MsgDefault{}
}

func CmdInvalid() tea.Msg {
	return MsgDefault{}
}

func ParseInput(msg string) tea.Cmd {
	for _, action := range actionsMap {
		for _, shorcut := range action.Shortcuts {
			if shorcut == msg {
				return action.Command
			}
		}
	}
	return CmdInvalid
}
