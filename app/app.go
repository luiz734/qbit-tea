package app

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dustin/go-humanize"
	"github.com/tubbebubbe/transmission"
)

var Timeout = time.Second * 2.0

func updateTableRows(m *Model, torrents transmission.Torrents) {
	rows := []table.Row{}
	for _, t := range torrents {
		strEta := formatTime(t.Eta)
		strPercentDone := fmt.Sprintf("%.1f", t.PercentDone*100.0)
		strStatus := TorrentStatus(t)
		strDown := humanize.Bytes(uint64(t.RateDownload))
		strUp := humanize.Bytes(uint64(t.RateUpload))
		rows = append(rows, table.Row{strEta, strPercentDone, strStatus, strDown, strUp, t.Name})
	}
	m.table.SetRows(rows)
}

type Model struct {
	client      *transmission.TransmissionClient
	address     string
	updateTimer timer.Model
	// torrents    transmission.Torrents
	table      table.Model
	windowSize struct {
		Width  int
		Height int
	}
	torrentRatio float64
}

func NewModel(updateTimer timer.Model, client *transmission.TransmissionClient, address string) Model {
	return Model{
		updateTimer: updateTimer,
		client:      client,
		table:       createTable([]table.Row{}, 0),
		address:     address,
	}
}

type msgStart struct{ cursor int }

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.updateTimer.Init(),
		CmdUpdate(m),
		tea.ClearScreen,
		func() tea.Msg { return msgStart{0} },
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
	case msgStart:
        log.Printf("before: %d", m.table.Cursor())
		// m.table.SetCursor(msg.cursor)
        log.Printf("after: %d", m.table.Cursor())
        return m, nil
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
		// Sort first, then update
		// m.torrents = msg.Torrents
		// m.torrents.SortByAddedDate(true)
		updateTableRows(&m, msg.Torrents)
		// var cmd tea.Cmd
		// m.table, cmd = m.table.Update(msg)
		// m.table = UpdateColumnsWidth(m.table, m.windowSize.Width)
		return m, nil

	case tea.WindowSizeMsg:
		m.windowSize = struct {
			Width  int
			Height int
		}{msg.Width - 2, msg.Height}
		// m.table = createTable(RenderTorrentTable(m.torrents, m.cursor), m.cursor)
		// m.table = UpdateColumnsWidth(m.table, m.windowSize.Width)
		// var tableCmd tea.Cmd
		// m.table, tableCmd = m.table.Update(m)
		// return m, tea.Batch(CmdUpdate(m), tea.ClearScreen, tableCmd)

	case tea.KeyMsg:
		switch msg.String() {
		// case "j", "down":
		// 	if m.cursor < 18 {
		// 		// m.cursor = m.cursor + 1
		// 		m.table.SetCursor(2)
		// 		log.Printf("cursor: %d", m.cursor)
		// 	}
		// case "k", "up":
		// 	if m.cursor > 0 {
		// 		// m.cursor = m.cursor - 1
		// 		m.table.SetCursor(2)
		// 		log.Printf("cursor: %d", m.cursor)
		// 	}
		// case "ctrl+c", "q":
		// 	return m, tea.Quit
		// case "u":
		// 	return m, CmdUpdate(m)
		// case "p":
		// 	return m, CmdToggle(m)
		// case "d":
		// 	return m, CmdRemove(m, false)
		// case "a":
		// 	s := NewDirModel(m)
		// 	return s.Update(nil)
		}
	}
	m.table, cmd = m.table.Update(msg)
	log.Printf("%d", m.table.Cursor())
	return m, cmd
}

func (m Model) View() string {
	var output strings.Builder

	// header := viewHeader(m)
	// output.WriteString(header)

	// output.WriteString("\n\n")
	output.WriteString(m.table.View())
	output.WriteString("\n\n")
	output.WriteString(m.table.HelpView())

	return output.String()
}
