package addtorrent

import (
	tea "github.com/charmbracelet/bubbletea"
)

type SubmitButton struct{}

func SubmitCmd() tea.Msg {
	return SubmitButton{}
}

type ModelButton struct {
	label   string
	focused bool
}

func NewModelButton(label string) *ModelButton {
	return &ModelButton{
		label: label,
	}
}

func (m *ModelButton) Init() tea.Cmd {
	return nil
}

func (m *ModelButton) Update(msg tea.Msg) (Focuser, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		magicNumber := 6
		paddingHorizontal := (msg.Width) / 2
		styleButton = styleButton.Margin(0, paddingHorizontal-magicNumber)
		styleButtonFocused = styleButtonFocused.Margin(0, paddingHorizontal-magicNumber)
	case tea.KeyMsg:
		space := " "
		if msg.String() == space {
			return m, SubmitCmd
		}
	}

	return m, nil
}

func (m *ModelButton) View() string {
	if m.focused {
		return styleButtonFocused.Render(m.label)
	}
	return styleButton.Render(m.label)
}

func (m *ModelButton) Focus() tea.Cmd {
	m.focused = true
	return nil
}

func (m *ModelButton) Blur() {
	m.focused = false
}

func (m *ModelButton) Value() string {
	return ""
}
