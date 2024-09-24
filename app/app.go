package app

import (
	"fmt"
	"log"
	"qbit-tea/input"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tubbebubbe/transmission"
)

var Timeout = time.Second * 2.0

type Model struct {
	client      *transmission.TransmissionClient
	updateTimer timer.Model
	torrents    transmission.Torrents
	cursor      int
	addMode     bool

	table table.Model
}

func NewModel(updateTimer timer.Model, client *transmission.TransmissionClient) Model {
	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Status", Width: 12},
		{Title: "Name", Width: 40},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	return Model{
		updateTimer: updateTimer,
		client:      client,
		table:       t,
	}
}

func (m Model) Init() tea.Cmd {
	log.Printf(fmt.Sprintf("%s", m.torrents))
	return tea.Batch(m.updateTimer.Init(), CmdUpdate(m))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case timer.TickMsg:
		var cmd tea.Cmd
		m.updateTimer, cmd = m.updateTimer.Update(msg)
		return m, cmd

	case timer.TimeoutMsg:
		m.updateTimer.Timeout = Timeout
		m.updateTimer.Init()
		return m, CmdUpdate(m)

	case MsgUpdate:
		m.torrents = msg.Torrents
		m.table.SetRows(RenderTorrentTable(m.torrents, m.cursor))
		return m, nil

	case input.MsgStart:
		return m, nil

	case MsgAdd:
		// m.torrents[m.cursor].Name += msg.Foo
		return m, nil

	case input.MsgMoveCursor:
		newPos := m.cursor + msg.Movement
		if newPos >= 0 && newPos <= len(m.torrents)-1 {
			m.cursor = newPos
		}
		m.table.SetCursor(m.cursor)
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
		case "a":
			return m, CmdAdd(m)
		default:
			return m, input.ParseInput(msg.String())
		}
	}

	return m, nil
}

func (m Model) View() string {
	var output strings.Builder
	timeoutSec := fmt.Sprintf("%d", (m.updateTimer.Timeout/1_000_000_000)+1)
	output.WriteString(fmt.Sprintf("Updating in %s...\n", timeoutSec))
	output.WriteString(m.table.View())
	return output.String()
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
