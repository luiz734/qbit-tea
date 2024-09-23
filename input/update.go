package input

import (
	"bufio"
	"bytes"
	"log"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

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

type MsgUpdate struct {
	Entries []Entry
}

func CmdUpdate() tea.Msg {
	args := []string{"--list"}
	cmd := exec.Command("transmission-remote", args...)
	stdout, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	return MsgUpdate{Entries: ParseEntries(stdout)}
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
        _=log.Default
		fileName := strings.Join(fields[8:], " ")
		_ = fileName
		entries = append(entries, NewEntry(fields[0], fileName, fields[7]))
	}

	return entries
}
