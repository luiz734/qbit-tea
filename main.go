package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
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
		entries = append(entries, NewEntry(fields[0], fields[8], fields[9]))
	}

	return entries
}

type model struct {
	Entries []Entry
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	var output strings.Builder
	for _, entry := range m.Entries {
		output.WriteString(fmt.Sprintf("%s %s %s", entry.Id, entry.Status, entry.Filename))
	}
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
