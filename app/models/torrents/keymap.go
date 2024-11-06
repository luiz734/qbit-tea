package torrents

import "github.com/charmbracelet/bubbles/key"

type Keymap struct {
	Up     key.Binding
	Down   key.Binding
	Update key.Binding
	Add    key.Binding
	Toggle key.Binding
	Delete key.Binding
	Info   key.Binding

	Quit key.Binding
	Help key.Binding
}

func (k Keymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Help, k.Quit}
}

func (k Keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Add},          // first column
		{k.Delete, k.Toggle, k.Update}, // second column
		{k.Info, k.Help, k.Quit},       // ...
	}
}

// DefaultKeyMap returns a set of pager-like default keybindings.
func DefaultKeymap() Keymap {
	return Keymap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Update: key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "update"),
		),
		Add: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add torrent"),
		),
		Toggle: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "start/stop"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete"),
		),
		Info: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "show info"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
	}
}
