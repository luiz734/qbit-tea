package torrents

import (
	"fmt"
	"log"
	"qbit-tea/app/models/addtorrent"
	"qbit-tea/app/models/errorscreen"
	"qbit-tea/app/models/torrentinfo"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tubbebubbe/transmission"
)

var Timeout = time.Second * 2.0

type windowSize struct {
	Width  int
	Height int
}

type Model struct {
	client  *transmission.TransmissionClient
	address string
	updateTimer timer.Model
	table       table.Model
	windowSize  windowSize
	help        help.Model
	keymap      Keymap
	// Selected row. Can be null before the data is loaded
	torrentsCount int
	torrent       *transmission.Torrent
}

func NewModel(updateTimer timer.Model, client *transmission.TransmissionClient, address string) Model {
	return Model{
		updateTimer: updateTimer,
		client:      client,
		table:       createTable([]table.Row{}, 0),
		address:     address,
		help:        help.New(),
		keymap:      DefaultKeymap(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.updateTimer.Init(),
		CmdUpdate(m),
		tea.ClearScreen,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.Help):
			m.help.ShowAll = !m.help.ShowAll
			// Return tea.ClearScreen fix a bug
			// that makes the help view not render
			// correctly after toggle
			return m, tea.ClearScreen
		// case key.Matches(msg, m.keymap.Up):
		// 	return m, tea.Quit
		// case key.Matches(msg, m.keymap.Down):
		// 	return m, tea.Quit
		case key.Matches(msg, m.keymap.Update):
			return m, CmdUpdate(m)
		case key.Matches(msg, m.keymap.Toggle):
			return m, CmdToggle(m)
		case key.Matches(msg, m.keymap.Delete):
			return m, CmdRemove(m, false)
		case key.Matches(msg, m.keymap.Delete):
			return m, CmdRemove(m, false)
		case key.Matches(msg, m.keymap.Add):
			s := addtorrent.InitialModel(m, m.windowSize.Width, m.windowSize.Height)
			return s, s.Init()
		case key.Matches(msg, m.keymap.Info):
			s := torrentinfo.InitialModel(m, m.windowSize.Width, m.windowSize.Height, m.torrent)
			return s, s.Init()
		}
	// Trigger after user select a dir and magnet
	case addtorrent.FormDataMsg:
		log.Printf("Got form data %+v", msg)
		if msg.Magnet == "" || msg.DownloadDir == "" {
			// User cancel the operation
			return m, m.updateTimer.Init()
		}
		log.Printf("target dir: %s\nMagnet: %s\n", msg.DownloadDir, msg.Magnet)
		addCommand, err := NewAddInDirCmdByMagnet(msg.Magnet, msg.DownloadDir)

		if err != nil {
			panic(err)
		}
		_, err = m.client.ExecuteCommand(addCommand)

		if err != nil {
			panic(err)
		}
		log.Printf("add torrent %s", msg)
		m.updateTimer = timer.NewWithInterval(Timeout, time.Millisecond)
		return m, m.updateTimer.Init()

	case timer.TickMsg:
		m.updateTimer, cmd = m.updateTimer.Update(msg)
		return m, cmd

	case timer.TimeoutMsg:
        log.Printf("Tick")
		m.updateTimer.Timeout = Timeout
		m.updateTimer.Init()
		return m, CmdUpdate(m)

	case MsgUpdate:
		updateTableRows(&m, msg.Torrents)
		return m, nil

	case MsgError:
		s := errorscreen.InitialModel(msg.prevModel, msg.title, msg.err, m.windowSize.Width, m.windowSize.Height)
		return s, s.Init()

	case tea.WindowSizeMsg:
		log.Print(msg)
		m.windowSize = windowSize{msg.Width, msg.Height}
		updateTableSize(&m.table, m.windowSize)
		m.help.Width = msg.Width
		return m, tea.Batch(CmdUpdate(m), tea.ClearScreen)
	}

	// Update styles
	sectionWidth := m.windowSize.Width / 3
	leftStyle = leftStyle.Width(sectionWidth)
	centerStyle = centerStyle.Width(sectionWidth)
	rightStyle = rightStyle.Width(sectionWidth)

	m.table, cmd = m.table.Update(msg)
	torrents, err := m.client.GetTorrents()
	m.torrentsCount = len(torrents)

	if err != nil {
		panic(err)
	}
	torrents.SortByAddedDate(true)

	// We dont use m.torrentsCount here because there
	// is some delay before it actually add the torrent
	// and we may get out of range
	// len(torrent) is always accurate
	if len(torrents) > 0 {
		m.torrent = &torrents[m.table.Cursor()]
	}

	return m, cmd
}

func (m Model) View() string {
	// Header
	headerView := viewHeader(m)
	// Torrents
	var torrentsView string
	if m.torrentsCount == 0 {
		torrentsView = "<add some torrents>"
	} else {
		torrentsView = m.table.View()
		// Adding a new line fix first
		// row of help not showing
		torrentsView += "\n"
	}
	// Help
	helpView := m.help.View(m.keymap)
	// Gap
	headerHeight := lipgloss.Height(headerView)
	torrentsHeight := lipgloss.Height(torrentsView)
	helpHeight := lipgloss.Height(helpView)
	magicNumber := 2
	gapBottom := m.windowSize.Height - headerHeight - torrentsHeight - helpHeight + magicNumber
	if gapBottom < 0 {
		gapBottom = 0
	}
	gapBottomView := strings.Repeat("\n", gapBottom)

	return fmt.Sprint(
		headerView,
		torrentsView,
		gapBottomView,
		helpView,
	)
}
