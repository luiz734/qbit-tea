package app

import (
	"fmt"
	"qbit-tea/util"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type inputModel struct {
	textInput    textinput.Model
	prevModel    tea.Model
	autocomplete string
}

func (m *inputModel) Init() tea.Cmd {
	return textinput.Blink
}

type inputMsg string

func NewInputModel(prevModel tea.Model) inputModel {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	// todo: fix hardcoded width
	ti.Width = 50

	return inputModel{
		textInput:    ti,
		prevModel:    prevModel,
	}
}

func (m *inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {
		// Allways handle quit first
		case tea.KeyCtrlC:
			return m, tea.Quit
		// Autocomplete the value from placeholder
		case tea.KeyTab:
			m.textInput.SetValue(m.autocomplete)
		// User sends the input back to the parent model
		case tea.KeyEnter:
			return m.prevModel.Update(inputMsg(m.textInput.Value()))
		// User cancel the operation
		case tea.KeyEscape:
			return m.prevModel.Update(inputMsg(""))
		}
	}

	// Automatically detect a magnet link on the clipboard and set as placeholder
    placeholder, err := clipboard.ReadAll()
    util.CheckError(err)
    m.textInput.Placeholder = "Paste the magnet link"
    m.autocomplete = ""
	if strings.HasPrefix(placeholder, "magnet:?xt=") {
		m.textInput.Placeholder = "Magnet detected. Press <Tab> to fill"
        m.autocomplete = placeholder
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m *inputModel) View() string {
	return fmt.Sprintf(
		"Add torrent\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to abort)",
	) + "\n"
}
