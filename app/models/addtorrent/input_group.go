package addtorrent

import (
	// "github.com/charmbracelet/bubbles/textinput"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

type TextInputFocuser struct {
	Model textinput.Model
}

func (t *TextInputFocuser) Update(msg tea.Msg) (Focuser, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Again, width must be smaller thant the window
		// This time, the width should fill the remaining space
		// Prompt takes some, also border
		// This 3 should be 2 (1 border on each side)
		widthPromp := len(t.Model.Prompt)
		magicNumber := widthPromp + 3
		t.Model.Width = msg.Width - magicNumber
	}
	updatedModel, cmd := t.Model.Update(msg)
	t.Model = updatedModel
	return t, cmd
}

func (t *TextInputFocuser) View() string {
	return t.Model.View()
}

func (t *TextInputFocuser) Focus() tea.Cmd {
	return t.Model.Focus()
}

func (t *TextInputFocuser) Blur() {
	t.Model.Blur()
}

func (t *TextInputFocuser) Value() string {
	return t.Model.Value()
}

// End wrapper

// Models that can receive/lose focus
type Focuser interface {
	Update(tea.Msg) (Focuser, tea.Cmd)
	View() string
	Focus() tea.Cmd
	Blur()
	Value() string
}

type InputGroup struct {
	inputGroup []Focuser
	focus      int
}

func (m InputGroup) Init() tea.Cmd {
	return nil
}

func (m InputGroup) Update(msg tea.Msg) (InputGroup, tea.Cmd) {
	var cmds []tea.Cmd

	for index, item := range m.inputGroup {
		var cmd tea.Cmd
		m.inputGroup[index], cmd = item.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m InputGroup) View() string {
	output := ""
	for index, item := range m.inputGroup {
		if index == m.focus {
			output += styleFocused.Render(item.View())
		} else {
			output += styleUnfocused.Render(item.View())
		}
	}
	return output
}

func NewInputGroup(f ...Focuser) InputGroup {
	var items []Focuser
	for _, item := range f {
		items = append(items, item)
	}
	if len(items) == 0 {
		log.Fatal("Can't create empty input group")
	}

	items[0].Focus()
	// Focus the first non empty
	// for index := range f {
	// 	if items[index].Value() == "" {
	// 		items[index].Focus()
	// 		break
	// 	}
	// }
	return InputGroup{
		inputGroup: items,
	}
}

func (g *InputGroup) SetFocus(index int) tea.Cmd {
	if index > len(g.inputGroup)-1 {
		return nil
	}
	g.focus = index
	for _, i := range g.inputGroup {
		i.Blur()
	}
	return g.inputGroup[g.focus].Focus()
}

func (g *InputGroup) FocusNext() tea.Cmd {
	if len(g.inputGroup) == 0 {
		log.Debug("Empty input group")
		return nil
	}
	g.focus += 1
	if g.focus > len(g.inputGroup)-1 {
		g.focus = 0
	}
	log.Debug("Update focused index", "index", g.focus)

	// range makes a copy
	// Use the index, not the copy
	for index := range g.inputGroup {
		g.inputGroup[index].Blur()
	}
	return g.inputGroup[g.focus].Focus()
}

func (g InputGroup) GetFocused() *Focuser {
	if len(g.inputGroup) == 0 {
		return nil
	}
	return &g.inputGroup[g.focus]
}

func (g InputGroup) GetFocusedIndex() int {
	return g.focus
}

func (g InputGroup) getFieldValue(index int) (string, error) {
	if index >= len(g.inputGroup) {
		return "", fmt.Errorf("out of range: len = %d, index = %d", len(g.inputGroup), index)
	}
	return g.inputGroup[index].Value(), nil
}

func (g InputGroup) GetFormData() (*FormDataMsg, error) {
	var err error
	var downloadDir, subDir, magnet string
	if downloadDir, err = g.getFieldValue(0); err != nil {
		panic(err)
	}
	if subDir, err = g.getFieldValue(1); err != nil {
		panic(err)
	}
	if magnet, err = g.getFieldValue(2); err != nil {
		panic(err)
	}
	if _, err := os.Stat(downloadDir); err != nil {
		return nil, fmt.Errorf("download dir %s doens't exists", downloadDir)
	}
	if !isMagnet(magnet) {
		return nil, fmt.Errorf("invalid magnet link \"%s\"", magnet)
	}

	// TODO: check for errors
	return &FormDataMsg{
		DownloadDir: filepath.Join(downloadDir, subDir),
		Magnet:      magnet,
	}, nil
}
