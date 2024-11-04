package errorscreen

import "github.com/charmbracelet/bubbles/key"

type keymap struct {
	Exit key.Binding
	Quit key.Binding
	Help key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.Exit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Exit, k.Help}, // first column
		{k.Quit},         // second column
	}
}

// DefaultKeyMap returns a set of pager-like default keybindings.
func DefaultKeyMap() keymap {
	return keymap{
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
		Exit: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("any", "go back"),
		),
	}
}
