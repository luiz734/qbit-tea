package errorscreen

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	errTitle string
	errDesc  string
}

func InitialModel(prevModel tea.Model, errTitle string, err error, width, height int) Model {
	return Model{
		prevModel: prevModel,
		help:      help.New(),
		keyMap:    DefaultKeyMap(),
		errTitle:  errTitle,
		errDesc:   err.Error(),
		width:     width,
		height:    height,
	}
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
		case key.Matches(msg, m.keyMap.Exit):
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, tea.ClearScreen
		default:
			if _, ok := m.prevModel.(quitModel); ok {
				return m, tea.Quit
			}
			return m.prevModel, nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Update the style based on window size
		// This magic numbers was found testing
		// For some reason, width them things not align properly
		magicNumber := 2
		styleError = styleError.Width(m.width - magicNumber)
		styleHelp = styleHelp.Width(m.width - magicNumber)

		// We comment out here to render new lines manually
		// Add spaces before and after later
		// styleOutput = styleOutput.Height(m.height - magicNumber)

		cmds = append(cmds, tea.ClearScreen)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	errorView := lipgloss.JoinVertical(lipgloss.Center,
		styleErrTitle.Render(m.errTitle),
		styleErrDesc.Render(m.errDesc),
	)
	errorView = styleError.Render(errorView)
	helpView := styleHelp.Render(m.help.View(m.keyMap))

	heightError := lipgloss.Height(errorView)
	heightHelp := lipgloss.Height(helpView)

	// Always center the error message
	gapTop := (m.height - heightError) / 2
	if gapTop < 0 {
		gapTop = 0
	}
	gapTopView := strings.Repeat("\n", gapTop)

	// If help has many lines, remove from bottom gap only
	gapBottom := gapTop - heightHelp + 1
	if gapBottom < 0 {
		gapBottom = 0
	}
	gapBottomView := strings.Repeat("\n", gapBottom)

	return fmt.Sprintf("%s%s%s%s",
		gapTopView,
		errorView,
		gapBottomView,
		helpView,
	)
}

func QuitModel() quitModel { return quitModel{} }

// Used to avoid pointer
// If the type of prevModel matches this type, quit the application
type quitModel struct{}

func (q quitModel) Init() tea.Cmd                       { return tea.Quit }
func (q quitModel) Update(tea.Msg) (tea.Model, tea.Cmd) { return q, func() tea.Msg { return tea.Quit } }
func (q quitModel) View() string                        { return "" }
