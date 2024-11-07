package torrentinfo

import (
	"fmt"
	"qbit-tea/app/models"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/muesli/reflow/wrap"
	"github.com/tubbebubbe/transmission"
)

type Model struct {
	// Go back to prev model
	prevModel tea.Model
	// help
	help help.Model
	// Keymaps
	keyMap keymap

	// Window size
	width  int
	height int

	// Error message
	torrent *transmission.Torrent
}

func InitialModel(prevModel tea.Model, width, height int, torrent *transmission.Torrent) Model {
	m := Model{
		prevModel: prevModel,
		width:     width,
		height:    height,
		help:      help.New(),
		keyMap:    DefaultKeyMap(),
		torrent:   torrent,
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.ClearScreen,
		func() tea.Msg {
			return tea.WindowSizeMsg{Width: m.width, Height: m.height}
		},
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var err error
	var cmds []tea.Cmd
	var cmd tea.Cmd
	_ = err
	_ = cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, tea.ClearScreen
		case key.Matches(msg, m.keyMap.Exit):
			return m.prevModel, nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		cmds = append(cmds, tea.ClearScreen)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	infoView := m.infoView()
	infoView = styleWrapper.Render(infoView) + "\n"
	helpView := styleHelp.Render(m.help.View(m.keyMap))

	infoHeight := lipgloss.Height(infoView)
	helpHeight := lipgloss.Height(helpView)

	// Try and error
	magicNumber := 1
	gapBottom := m.height - infoHeight - helpHeight + magicNumber
	if gapBottom < 0 {
		gapBottom = 0
	}
	gapBottomView := strings.Repeat("\n", gapBottom)

	return fmt.Sprintf("%s%s%s",
		infoView,
		gapBottomView,
		helpView,
	)
}

// TODO: handle small width better
func (m Model) infoView() string {
	val := reflect.ValueOf(m.torrent).Elem()
	var output strings.Builder
	// -2 because last 2 fields hold error values
	for i := 0; i < val.NumField()-2; i++ {
		typeName := val.Type().Field(i).Name
		value := fmt.Sprint(val.Field(i))
		var err error
		if value, err = humanizeValue(typeName, value); err != nil {
			errMsg := fmt.Sprintf("humanize value: %v", err)
			panic(errMsg)
		}

		labelView := styleLabel.Render(fmt.Sprintf("%s", typeName))
		valueView := styleValue.Render(fmt.Sprint(value))
		labelWidth := lipgloss.Width(labelView)
		valueWidth := lipgloss.Width(valueView)
		// Border and padding add some extra width
		// TODO: calculate using lipgloss

		// padding + border
		magicNumber := 6 + 2
		nDots := m.width - labelWidth - valueWidth - magicNumber
		if nDots < 0 {
			nDots = 0
		}
		separatorView := styleSeparator.Render(strings.Repeat(".", nDots))
		if nDots != 0 {
			output.WriteString(lipgloss.JoinHorizontal(
				lipgloss.Left,
				labelView,
				separatorView,
				valueView,
			))
		} else {
			output.WriteString(labelView)
			output.WriteString(valueView)
		}
		// Don't add newline for the last field
		if i < val.NumField()-3 {
			output.WriteRune('\n')
		}
	}
	// Same as magicNumber
	return wrap.String(output.String(), m.width-8)
}

// Format SOME values
// We take the value as string because it's simpler
// PERF: work with reflect library properly
func humanizeValue(fieldName, value string) (string, error) {
	switch fieldName {
	case "AddedDate":
		asInt64, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return "", fmt.Errorf("parsing AddedDate: %w", err)
		}
		return humanize.Time(time.Unix(asInt64, 0)), nil
	case "Status":
		asInt64, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return "", fmt.Errorf("parsing Status: %w", err)
		}
		// HACK: this function works with torrents only
		// Maybe use an interface instead
		return models.TorrentStatus(transmission.Torrent{Status: int(asInt64)}), nil
	case "PercentDone":
		asFloat, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return "", fmt.Errorf("parsing PercentDone: %w", err)
		}
		return fmt.Sprintf("%s%%", humanize.Ftoa(asFloat*100.0)), nil
	}
	return value, nil
}
