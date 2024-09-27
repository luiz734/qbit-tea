package app

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type downloadSubDir struct {
	textInput textinput.Model
	prevModel tea.Model
}

func (m *downloadSubDir) Init() tea.Cmd {
	return textinput.Blink
}

type downloadSubDirMsg string

func NewDownloadSubDirModel(prevModel tea.Model) downloadSubDir {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	// todo: fix hardcoded width
	ti.Width = 50

	return downloadSubDir{
		textInput: ti,
		prevModel: prevModel,
	}
}

func (m *downloadSubDir) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {
		// Allways handle quit first
		case tea.KeyCtrlC:
			return m, tea.Quit
		// User sends the input back to the parent model
		case tea.KeyEnter:
			return m.prevModel.Update(downloadSubDirMsg(m.textInput.Value()))
		// User cancel the operation
		case tea.KeyEscape:
			return m.prevModel.Update(downloadSubDirMsg(""))
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m *downloadSubDir) View() string {
	return fmt.Sprintf(
		"Subdir\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to abort)",
	) + "\n"
}
