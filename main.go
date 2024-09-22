package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"qbit-tea/input"
	"qbit-tea/transmission"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

type Entry struct {
	Id       string
	Filename string
	Status   string
}

func NewEntry(id string, filename string, status string) Entry {
	return Entry{
		Id:       id,
		Filename: filename,
		Status:   status,
	}
}

func ParseEntries(output []byte) []Entry {

	scanner := bufio.NewScanner(bytes.NewReader(output))

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// remove header and "sum" line
	if len(lines) > 2 {
		lines = lines[1 : len(lines)-1]
	}

	var entries []Entry
	for _, line := range lines {
		fields := strings.Fields(line)
		fileName := strings.Join(fields[9:], " ")
		_ = fileName
		entries = append(entries, NewEntry(fields[0], fileName, fields[8]))
	}

	return entries
}

type model struct {
	Entries []Entry
	cursor  int
	addMode bool
}

func (m model) Init() tea.Cmd {
	return nil
}

type actionMsg struct {
	helpItem input.UserAction
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case input.MsgUpdate:
		// m.Entries[m.cursor].Filename += msg.Output
		output, err := transmission.TransmissionList()
		checkError(err)
		m.Entries = ParseEntries(output)
		return m, nil

    case input.MsgAdd:
		m.Entries[m.cursor].Filename += msg.Foo
        return m, nil


	case input.MsgMoveCursor:
		newPos := m.cursor + msg.Movement
		if newPos >= 0 && newPos <= len(m.Entries)-1 {
			m.cursor = newPos
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		default:
			return m, input.ParseInput(msg.String())
		}
	}

	return m, nil
}

func (m model) View() string {
	var output strings.Builder
	for index, entry := range m.Entries {
		if index == m.cursor {
			output.WriteRune('>')
		} else {
			output.WriteRune(' ')
		}
		output.WriteString(fmt.Sprintf("%-4s %-10s %s\n", entry.Id, entry.Status, entry.Filename))
	}
	output.WriteString(input.HelpMsg())
	return output.String()
}

func main() {
	output, err := transmission.TransmissionList()
	checkError(err)

	if _, err := tea.NewProgram(model{Entries: ParseEntries(output)}).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}

}
