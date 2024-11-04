package addtorrent

import "github.com/charmbracelet/bubbles/key"

type AddTorrentKeymap struct {
	Next  key.Binding
	Add   key.Binding
	Clear key.Binding
	Quit  key.Binding
	Abort key.Binding
	Help  key.Binding
}

func (k AddTorrentKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Next, k.Add, k.Abort, k.Help}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k AddTorrentKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Add, k.Next, k.Clear},  // first column
		{k.Quit, k.Abort, k.Help}, // second column
	}
}

// DefaultKeyMap returns a set of pager-like default keybindings.
func DefaultAddTorrentKeyMap() AddTorrentKeymap {
	return AddTorrentKeymap{
		Next: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next field"),
		),
		Add: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "add torrent"),
		),
		Clear: key.NewBinding(
			key.WithKeys("C-l"),
			key.WithHelp("clear", "clear filed"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("quit", "quit application"),
		),
		Abort: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
	}
}
