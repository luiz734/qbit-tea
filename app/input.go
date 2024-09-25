package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type inputModel struct {
	textInput textinput.Model
	prevModel tea.Model
}

func (m *inputModel) Init() tea.Cmd {
	return textinput.Blink
}

type inputMsg string

func NewInputModel(prevModel tea.Model) inputModel {
	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return inputModel{
		textInput: ti,
        prevModel: prevModel,
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
		// User sends the input back to the parent model
		case tea.KeyEnter:
			return m.prevModel.Update(inputMsg(m.textInput.Value()))
		// User cancel the operation
		case tea.KeyEscape:
			return m, nil
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m *inputModel) View() string {
	return fmt.Sprintf(
		"What’s your favorite Pokémon?\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}
