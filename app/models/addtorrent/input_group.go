package addtorrent

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

type InputGroup struct {
	inputGroup []textinput.Model
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
	for _, item := range m.inputGroup {
		output += styleMagnet.Render(item.View())
	}
	return output
}

func NewInputGroup(f ...textinput.Model) InputGroup {
	var items []textinput.Model
	for _, item := range f {
		items = append(items, item)
	}
	if len(items) == 0 {
		log.Fatal("Can't create empty input group")
	}
	items[0].Focus()
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
	for index, _ := range g.inputGroup {
		g.inputGroup[index].Blur()
	}
	return g.inputGroup[g.focus].Focus()
}

func (g InputGroup) GetFocused() *textinput.Model {
	if len(g.inputGroup) == 0 {
		return nil
	}
	return &g.inputGroup[g.focus]
}

func (g InputGroup) GetFocusedIndex() int {
	return g.focus
}
