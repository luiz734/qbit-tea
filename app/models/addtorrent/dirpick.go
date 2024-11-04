package addtorrent

import (
	"fmt"
	"io"
	"qbit-tea/config"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	funk "github.com/thoas/go-funk"
)

type dirPickModel struct {
	list    list.Model
	focused bool
}

func (m dirPickModel) Init() tea.Cmd {
	return nil
}

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}
	// str := fmt.Sprintf("%d. %s", index+1, i)
	str := fmt.Sprintf("%s", i)
	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}


func NewDickPickModel() *dirPickModel {
	// Convert []string to []list.Item
	// dirs := []string{JellyShowsDir, JellyMoviesDir}
	dirs := config.Cfg.DownloadDirs
	// sort.Strings(dirs)
	items := funk.Map(dirs, func(s string) list.Item {
		return item(s)
	}).([]list.Item)

	l := list.New(items, itemDelegate{}, 10, 8)
	l.SetFilteringEnabled(false)
	l.SetShowTitle(true)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.Title = "Download directory"
    l.Styles.Title = titleStyle
    l.SetWidth(len(l.Title) + 10)

	return &dirPickModel{
		list: l,
	}
}

type inputMsg string

func (m *dirPickModel) Update(msg tea.Msg) (Focuser, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
        // This works fine
        // Changing Styles.Title does nothing to the width
        magicNumber := 2
        m.list.Styles.TitleBar = m.list.Styles.TitleBar.Width(msg.Width - magicNumber)

		// m.list.SetWidth(msg.Width)
		_ = msg
		return m, nil
	}
	if !m.focused {
		return m, nil
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *dirPickModel) View() string {
	return m.list.View()
}

func (m *dirPickModel) GetPath() string {
	item, ok := m.list.SelectedItem().(item)
	if ok {
		return string(item)
	}
	return ""
}

func (m *dirPickModel) Focus() tea.Cmd {
	m.focused = true
	return nil
}

func (m *dirPickModel) Blur() {
	m.focused = false
}

func (m dirPickModel) Value() string {
	item, ok := m.list.SelectedItem().(item)
	if ok {
		return string(item)
	}
    log.Warnf("Value is empty")
	return ""
}
