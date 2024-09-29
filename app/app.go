package app

import (
	"log"
	"qbit-tea/util"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tubbebubbe/transmission"
)

var Timeout = time.Second * 2.0

type windowSize struct {
	Width  int
	Height int
}

type Model struct {
	client      *transmission.TransmissionClient
	address     string
	updateTimer timer.Model
	// torrents    transmission.Torrents
	table        table.Model
	windowSize   windowSize
    // Selected row. Can be null before the data is loaded
	torrent      *transmission.Torrent
}

func NewModel(updateTimer timer.Model, client *transmission.TransmissionClient, address string) Model {
	return Model{
		updateTimer: updateTimer,
		client:      client,
		table:       createTable([]table.Row{}, 0),
		address:     address,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.updateTimer.Init(),
		CmdUpdate(m),
		tea.ClearScreen,
	)
}

// func NewAddInDirCmdByMagnet(magnetLink string, path string) (*transmission.Command, error) {
// 	cmd, _ := transmission.NewAddCmd()
// 	cmd.Arguments.Filename = magnetLink
// 	// Can't check if it's a dir on remote hosts
// 	cmd.Arguments.DownloadDir = path
// 	return cmd, nil
// }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// User should always be able to quit
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	switch msg := msg.(type) {
	// Trigger after user select a dir and magnet
	// case dirMsg:
	// 	if msg.magnet == "" || msg.downloadDir == "" {
	// 		// User cancel the operation
	// 		return m, m.updateTimer.Init()
	// 	}
	// 	log.Printf("Target dir: %s\nMagnet: %s\n", msg.downloadDir, msg.magnet)
	// 	addCommand, err := NewAddInDirCmdByMagnet(msg.magnet, msg.downloadDir)
	// 	util.CheckError(err)
	// 	_, err = m.client.ExecuteCommand(addCommand)
	// 	util.CheckError(err)
	// 	log.Printf("Add torrent %s", msg)
	// 	m.updateTimer = timer.NewWithInterval(Timeout, time.Millisecond)
	// 	return m, m.updateTimer.Init()

	case timer.TickMsg:
		m.updateTimer, cmd = m.updateTimer.Update(msg)
		return m, cmd

	case timer.TimeoutMsg:
		m.updateTimer.Timeout = Timeout
		m.updateTimer.Init()
		return m, CmdUpdate(m)

	case MsgUpdate:
		updateTableRows(&m, msg.Torrents)
		return m, nil

	case tea.WindowSizeMsg:
		m.windowSize = windowSize{msg.Width, msg.Height}
		updateTableSize(&m.table, m.windowSize)
		return m, tea.Batch(CmdUpdate(m), tea.ClearScreen)

	case tea.KeyMsg:
		switch msg.String() {
		// case "ctrl+c", "q":
		// 	return m, tea.Quit
		// case "u":
		// 	return m, CmdUpdate(m)
		case "p":
			return m, CmdToggle(m)
		// case "d":
		// 	return m, CmdRemove(m, false)
		case "a":
			s := NewDirModel(m)
			return s.Update(nil)
		}
	}
	m.table, cmd = m.table.Update(msg)
	torrents, err := m.client.GetTorrents()
	util.CheckError(err)
	torrents.SortByAddedDate(true)
	m.torrent = &torrents[m.table.Cursor()]

	log.Printf("%d", m.table.Cursor())
	return m, cmd
}

func (m Model) View() string {
	var output strings.Builder

	header := viewHeader(m)
	output.WriteString(header)
	output.WriteString("\n\n")

	output.WriteString(m.table.View())
	output.WriteString("\n\n")
	output.WriteString(m.table.HelpView())

	return output.String()
}
