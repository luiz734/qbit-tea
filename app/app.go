package app

import (
	"fmt"
	"log"
	"qbit-tea/input"
	"qbit-tea/util"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tubbebubbe/transmission"
)

var Timeout = time.Second * 2.0

type Model struct {
	client      *transmission.TransmissionClient
	address     string
	updateTimer timer.Model
	torrents    transmission.Torrents
	cursor      int
	table       table.Model
	windowSize  struct {
		Width  int
		Height int
	}
	torrentRatio float64
}

func NewModel(updateTimer timer.Model, client *transmission.TransmissionClient, address string) Model {
	columns := []table.Column{
		{Title: "ETA", Width: 5},
		{Title: "%", Width: 6},
		{Title: "Status", Width: 12},
		{Title: "Down", Width: 8},
		{Title: "Up", Width: 8},
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
		address:     address,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.updateTimer.Init(), CmdUpdate(m), tea.ClearScreen)
}

func NewAddInDirCmdByMagnet(magnetLink string, path string) (*transmission.Command, error) {
	cmd, _ := transmission.NewAddCmd()
	cmd.Arguments.Filename = magnetLink
	// Can't check if it's a dir on remote hosts
	cmd.Arguments.DownloadDir = path
	return cmd, nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// User should always be able to quit
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	switch msg := msg.(type) {

	// Trigger after user select a dir and magnet
	case dirMsg:
		if msg.magnet == "" || msg.downloadDir == "" {
			// User cancel the operation
			return m, m.updateTimer.Init()
		}
		log.Printf("Target dir: %s\nMagnet: %s\n", msg.downloadDir, msg.magnet)
		addCommand, err := NewAddInDirCmdByMagnet(msg.magnet, msg.downloadDir)
		util.CheckError(err)
		_, err = m.client.ExecuteCommand(addCommand)
		util.CheckError(err)
		log.Printf("Add torrent %s", msg)
		m.updateTimer = timer.NewWithInterval(Timeout, time.Millisecond)
		return m, m.updateTimer.Init()

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
		m.torrents.SortByAddedDate(true)
		m.table.SetRows(RenderTorrentTable(m.torrents, m.cursor))
		m.table = UpdateColumnsWidth(m.table, m.windowSize.Width)
		return m, nil

	case input.MsgMoveCursor:
		newPos := m.cursor + msg.Movement
		if newPos >= 0 && newPos <= len(m.torrents)-1 {
			m.cursor = newPos
		}
		torrent := m.torrents[m.cursor]
		m.torrentRatio = torrent.UploadRatio
		m.table.SetCursor(m.cursor)
		return m, nil

	case tea.WindowSizeMsg:
		m.windowSize = struct {
			Width  int
			Height int
		}{msg.Width - 2, msg.Height}
		m.table = UpdateColumnsWidth(m.table, m.windowSize.Width)
		m.table = table.New(
			table.WithColumns(m.table.Columns()),
			table.WithFocused(true),
			table.WithHeight(msg.Height-2),
		)

		return m, tea.Batch(CmdUpdate(m), tea.ClearScreen)

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
			s := NewDirModel(m)
			return s.Update(nil)
			// s := NewInputModel(m)
			// return s.Update(nil)
		default:
			return m, input.ParseInput(msg.String())
		}
	}
	return m, nil
}

func (m Model) View() string {
	var output strings.Builder

	header := viewHeader(m)
	output.WriteString(header)

	output.WriteString("\n\n")
	output.WriteString(m.table.View())
	return output.String()
}

func viewHeader(m Model) string {
	sectionWidth := m.windowSize.Width / 3
	leftStyle := lipgloss.NewStyle().
		Width(sectionWidth).Align(lipgloss.Left).MarginLeft(1)
	centerStyle := lipgloss.NewStyle().
		Width(sectionWidth).Align(lipgloss.Center).Margin(0)
	rightStyle := lipgloss.NewStyle().
		Width(sectionWidth).Align(lipgloss.Right).MarginRight(1)

	// Left
	strRatio := fmt.Sprintf("%0.1f", m.torrentRatio)
	textLeft := leftStyle.Render(strRatio)

	// Center
	var styleAddress = lipgloss.NewStyle().
		Align(lipgloss.Center).Bold(true)
	address := styleAddress.Render(fmt.Sprintf("%s", m.address))
	textCenter := centerStyle.Render(address)

	// Right
	r := ":D"
	if m.torrentRatio < 1.0 {
		r = ":("
	} else if m.torrentRatio < 2.0 {
		r = ":|"
	}
	textRight := rightStyle.Render(r)

	return lipgloss.JoinHorizontal(lipgloss.Top, textLeft, textCenter, textRight)
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
