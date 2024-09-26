package app

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	funk "github.com/thoas/go-funk"
)

type dirModel struct {
	list        list.Model
	downloadDir string
	prevModel   tea.Model
}

func (m *dirModel) Init() tea.Cmd {
	return nil
}

type dirMsg struct {
	downloadDir string
	magnet      string
}

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

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

func NewDirModel(prevModel tea.Model, dirs []string) dirModel {
	// Convert []string to []list.Item
	items := funk.Map(dirs, func(s string) list.Item {
		return item(s)
	}).([]list.Item)

	l := list.New(items, itemDelegate{}, 20, 20)
	l.SetFilteringEnabled(false)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	_ = l

	return dirModel{
		list:      l,
		prevModel: prevModel,
	}
}

func (m *dirModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case inputMsg:
		return m.prevModel.Update(dirMsg{
			downloadDir: m.downloadDir,
			magnet:      string(msg),
		})
	case tea.KeyMsg:
		switch msg.Type {
		// Allways handle quit first
		case tea.KeyCtrlC:
			return m, tea.Quit
		// User sends the input back to the parent model
		case tea.KeyEnter:
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.downloadDir = string(i)
				s := NewInputModel(m)
				return s.Update(nil)
			}
		// User cancel the operation
		case tea.KeyEscape:
			return m.prevModel.Update(dirMsg{})
		}
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *dirModel) View() string {
	return "\n" + m.list.View()
}
