package input

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type UserAction struct {
	Shortcut    string
	Description string
	Command     tea.Cmd
}

var actionsMap = []UserAction{
	{"a", "Add", CmdAdd},
	{"u", "Update", CmdUpdate},
	{"k", "Move Up", CmdMoveCursorUp},
	{"j", "Move Down", CmdMoveCursorDown},
	{"p", "Pause/Unpause", CmdDefault},
	{"d", "Delete", CmdDefault},
	{"x", "Delete ad emove files", CmdDefault},
	{"q", "Quit", tea.Quit},
}

func HelpMsg() string {

	var output strings.Builder
	output.WriteString("---\n")
	for _, d := range actionsMap {
		output.WriteString(fmt.Sprintf("%-2s -> %s\n", d.Shortcut, d.Description))
	}
	return output.String()
}

type MsgUpdate struct {
	Output string
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

func CmdUpdate() tea.Msg {
	return MsgUpdate{"done update"}
}

func CmdAdd() tea.Msg {
	return MsgAdd{"added"}
}

func CmdDefault() tea.Msg {
	return MsgDefault{}
}

func CmdInvalid() tea.Msg {
	return MsgDefault{}
}

func ParseInput(msg string) tea.Cmd {
	for _, a := range actionsMap {
		if a.Shortcut == msg {
			return a.Command
		}
	}
	return CmdInvalid
}
