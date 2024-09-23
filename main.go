package main

import (
	"fmt"
	"os"
	"qbit-tea/input"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

var timeout = time.Second * 5.0

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

type model struct {
	UpdateTimer timer.Model
	Entries     []input.Entry
	cursor      int
	addMode     bool
}

func (m model) Init() tea.Cmd {
	return m.UpdateTimer.Init()
}

type actionMsg struct {
	helpItem input.UserAction
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case timer.TickMsg:
		var cmd tea.Cmd
		m.UpdateTimer, cmd = m.UpdateTimer.Update(msg)
		return m, cmd

	case timer.TimeoutMsg:
        m.UpdateTimer.Timeout = timeout
		m.UpdateTimer.Init()
		return m, input.CmdUpdate

	case input.MsgUpdate:
		m.Entries = msg.Entries
		return m, nil

	case input.MsgAdd:
		m.Entries[m.cursor].Filename += msg.Foo
		return m, nil

	case input.MsgMoveCursor:
		newPos := m.cursor + msg.Movement
		if newPos >= 0 && newPos <= len(m.Entries)-1 {
			m.cursor = newPos
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		default:
			return m, input.ParseInput(msg.String())
		}
	}

	return m, nil
}

func (m model) View() string {
	var output strings.Builder
	timeoutSec := fmt.Sprintf("%d", (m.UpdateTimer.Timeout / 1_000_000_000)+1)
	output.WriteString(fmt.Sprintf("Updating in %s...\n", timeoutSec))
	for index, entry := range m.Entries {
		if index == m.cursor {
			output.WriteRune('>')
		} else {
			output.WriteRune(' ')
		}
		output.WriteString(fmt.Sprintf("%-4s %-10s %s\n", entry.Id, entry.Status, entry.Filename))
	}
	output.WriteString(input.HelpMsg())
	return output.String()
}

func main() {
	// output, err := transmission.TransmissionList()
	// checkError(err)

	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	if _, err := tea.NewProgram(model{
        UpdateTimer: timer.NewWithInterval(timeout, time.Millisecond),
    }).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}

}
