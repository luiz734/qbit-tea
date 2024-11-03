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

// var JellyShowsDir = "/jellyfin/series"
// var JellyMoviesDir = "/jellyfin/movies"

type dirPickModel struct {
	list list.Model
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

func NewDickPickModel() dirPickModel {
	// Convert []string to []list.Item
	// dirs := []string{JellyShowsDir, JellyMoviesDir}
	dirs := append(config.Cfg.MoviesDirs, config.Cfg.ShowsDirs...)
    sort.Strings(dirs)
	items := funk.Map(dirs, func(s string) list.Item {
		return item(s)
	}).([]list.Item)

	l := list.New(items, itemDelegate{}, 5, 6)
	l.SetFilteringEnabled(false)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)

	return dirPickModel{
		list: l,
	}
}

type inputMsg string

func (m dirPickModel) Update(msg tea.Msg) (dirPickModel, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m dirPickModel) View() string {
	return m.list.View()
}

func (m dirPickModel) GetPath() string {
	item, ok := m.list.SelectedItem().(item)
	if ok {
		return string(item)
	}
	return ""
}
