package main

import (
	"fmt"
	"log"
	"os"
	"qbit-tea/input"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/tubbebubbe/transmission"
	// "github.com/pyed/transmission"
)

var timeout = time.Second * 2.0

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func CmdUpdate(m model) tea.Cmd {
	return func() tea.Msg {
		torrents, err := m.client.GetTorrents()
		checkError(err)
		return MsgUpdate{Torrents: torrents}
	}
}

func CmdToggle(m model) tea.Cmd {
	return func() tea.Msg {
		torrent := m.Torrents[m.cursor]
		switch torrent.Status {
		case StatusStopped:
			m.client.StartTorrent(torrent.ID)
		case StatusDownloading, StatusSeeding:
			m.client.StopTorrent(torrent.ID)
		}
		return MsgToggle{}
	}
}

func CmdRemove(m model, deleteData bool) tea.Cmd {
	return func() tea.Msg {
		torrent := m.Torrents[m.cursor]
        deleteCommand, err := transmission.NewDelCmd(torrent.ID, deleteData)
        output, err :=m.client.ExecuteCommand(deleteCommand)
        log.Println(fmt.Sprintf("%s", output))
        checkError(err)

        return nil
	}
}

type MsgToggle struct{}

type model struct {
	client      *transmission.TransmissionClient
	UpdateTimer timer.Model
	Torrents    transmission.Torrents
	cursor      int
	addMode     bool
}

func (m model) Init() tea.Cmd {
	log.Printf(fmt.Sprintf("%s", m.Torrents))
	return m.UpdateTimer.Init()
}

type actionMsg struct {
	helpItem input.UserAction
}

type MsgUpdate struct{ Torrents transmission.Torrents }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case timer.TickMsg:
		var cmd tea.Cmd
		m.UpdateTimer, cmd = m.UpdateTimer.Update(msg)
		return m, cmd

	case timer.TimeoutMsg:
		m.UpdateTimer.Timeout = timeout
		m.UpdateTimer.Init()
		return m, CmdUpdate(m)

	case MsgUpdate:
		m.Torrents = msg.Torrents
		return m, nil

	case input.MsgStart:
		return m, nil

	case input.MsgAdd:
		m.Torrents[m.cursor].Name += msg.Foo
		return m, nil

	case input.MsgMoveCursor:
		newPos := m.cursor + msg.Movement
		if newPos >= 0 && newPos <= len(m.Torrents)-1 {
			m.cursor = newPos
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "u":
			return m, CmdUpdate(m)
		case "p":
			return m, CmdToggle(m)
        case "d":
            return m, CmdRemove(m, false)
		default:
			return m, input.ParseInput(msg.String())
		}
	}

	return m, nil
}

func (m model) View() string {
	var output strings.Builder
	timeoutSec := fmt.Sprintf("%d", (m.UpdateTimer.Timeout/1_000_000_000)+1)
	output.WriteString(fmt.Sprintf("Updating in %s...\n", timeoutSec))
	for index, entry := range m.Torrents {
		if index == m.cursor {
			output.WriteRune('>')
		} else {
			output.WriteRune(' ')
		}
		output.WriteString(fmt.Sprintf("%-4d %-10s %s\n", entry.ID, TorrentStatus(entry), entry.Name))
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
	client := transmission.New("http://127.0.0.1:9091", "", "")
	// checkError(err)
	_, err := client.GetTorrents()
	if err != nil {
		log.Panic(err)
	}

	if _, err := tea.NewProgram(model{
		UpdateTimer: timer.NewWithInterval(timeout, time.Millisecond),
		client:      &client,
	}).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}

}

const (
	StatusStopped = iota
	StatusCheckPending
	StatusChecking
	StatusDownloadPending
	StatusDownloading
	StatusSeedPending
	StatusSeeding
)

func TorrentStatus(status transmission.Torrent) string {
	switch status.Status {
	case StatusStopped:
		return "Stopped"
	case StatusCheckPending:
		return "Check waiting"
	case StatusChecking:
		return "Checking"
	case StatusDownloadPending:
		return "Download waiting"
	case StatusDownloading:
		return "Downloading"
	case StatusSeedPending:
		return "Seed waiting"
	case StatusSeeding:
		return "Seeding"
	default:
		return "unknown"
	}
}
