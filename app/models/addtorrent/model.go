package addtorrent

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	// Go back to prev model
	prevModel *tea.Model
	// help
	help help.Model
	// Keymaps
	keyMap AddTorrentKeymap

	// Window size
	width  int
	height int
	// Things can can be focused
	inputs InputGroup
}

func InitialModel(prevModel *tea.Model) Model {
	ti := textinput.New()
    // TOD o: placeholder + margin add new lines
	ti.Placeholder = "paste the magnet link"
    ti.Prompt = "Magnet: "
	ti.CharLimit = 156
	ti.PromptStyle = stylePrompt
	ti.SetValue(GetMagnetFromClipboard())
    ti.Width = 5
	// ti.PlaceholderStyle = stylePlaceholder

	sd := textinput.New()
	sd.Placeholder = "leave empty to skip"
	sd.Prompt = "Subdir: "
	sd.CharLimit = 156
	sd.PromptStyle = stylePrompt
	// sd.PlaceholderStyle = stylePlaceholder
    sd.Width = 5

	return Model{
		prevModel: prevModel,
		help:      help.New(),
		keyMap:    DefaultAddTorrentKeyMap(),
		inputs: NewInputGroup(
			NewDickPickModel(),
			&TextInputFocuser{sd},
			&TextInputFocuser{ti},
		),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.ClearScreen, textinput.Blink)
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
			// Warning: User can't type ? anymore
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		case key.Matches(msg, m.keyMap.Next):
			return m, m.inputs.FocusNext()
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		magicNumber := 2
		styleHelp = styleHelp.Width(m.width - magicNumber)
        // Don't update the width here
        // The width depends on the prompt
        // Update in input_group.go instead 
		// styleUnfocused = styleUnfocused.Width(m.width - magicNumber)
		// styleFocused = styleFocused.Width(m.width - magicNumber)

		cmds = append(cmds, tea.ClearScreen)
	}

	m.inputs, cmd = m.inputs.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	formView := lipgloss.JoinVertical(lipgloss.Left,
		m.inputs.View(),
	)
	helpView := styleHelp.Render(m.help.View(m.keyMap))

	heightForm := lipgloss.Height(formView)
	heightHelp := lipgloss.Height(helpView)

	// Always center the error message
	// gapTop := (m.height - heightForm) / 2
	// gapTop := 0
	// if gapTop < 0 {
	// 	gapTop = 0
	// }
	// gapTopView := strings.Repeat("\n", gapTop)

	// If help has many lines, remove from bottom gap only
	gapBottom := m.height - heightHelp - heightForm + 1
	if gapBottom < 0 {
		gapBottom = 0
	}
	gapBottomView := strings.Repeat("\n", gapBottom)

	return fmt.Sprintf("%s%s%s",
		formView,
		gapBottomView,
		helpView,
	)
}

func GetMagnetFromClipboard() string {
	placeholder, err := clipboard.ReadAll()
	// Clipboard empty
	if err != nil {
		placeholder = ""
	}
	// if strings.HasPrefix(placeholder, "magnet:?xt=") {
	if strings.HasPrefix(placeholder, "foo") {
		return placeholder
	}
	return ""
}
