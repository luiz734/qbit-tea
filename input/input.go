package input

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type UserAction struct {
	Shortcuts   []string
	Description string
	Command     tea.Cmd
}

var actionsMap = []UserAction{
	// {[]string{"a"}, "Add", CmdAdd},
	// {[]string{"u"}, "Update", CmdUpdate},
	{[]string{"k", "up"}, "Move Up", CmdMoveCursorUp},
	{[]string{"j", "down"}, "Move Down", CmdMoveCursorDown},
	// {[]string{"p"}, "Pause/Unpause", CmdStart()},
	// {[]string{"d"}, "Delete", CmdDefault},
	{[]string{"x"}, "Delete and remove files", CmdDefault},
	// {[]string{"q"}, "Quit", tea.Quit},
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

func CmdMoveCursorUp() tea.Msg {
	return MsgMoveCursor{-1}
}
func CmdMoveCursorDown() tea.Msg {
	return MsgMoveCursor{1}
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
