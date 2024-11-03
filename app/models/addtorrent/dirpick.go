package addtorrent

import (
	"fmt"
	"io"
	"qbit-tea/config"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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

func isMovieDir(path string) bool {
	for _, s := range config.Cfg.MoviesDirs {
		if s == path {
			return true
		}
	}
	return false
}

func NewDickPickModel() *dirPickModel {
	// Convert []string to []list.Item
	// dirs := []string{JellyShowsDir, JellyMoviesDir}
	dirs := append(config.Cfg.MoviesDirs, config.Cfg.ShowsDirs...)
	sort.Strings(dirs)
	items := funk.Map(dirs, func(s string) list.Item {
		return item(s)
	}).([]list.Item)

	l := list.New(items, itemDelegate{}, 10, 8)
	l.SetFilteringEnabled(false)
	l.SetShowTitle(true)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.Styles.Title = titleStyle
	l.Title = "Download direcotry"

	return &dirPickModel{
		list: l,
	}
}

type inputMsg string

func (m *dirPickModel) Update(msg tea.Msg) (Focuser, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
        // This affects the title, but doesn't make the model bigger
		m.list.SetWidth(msg.Width)
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
	return ""
}
