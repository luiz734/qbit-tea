package input
//
// import (
// 	"bufio"
// 	"bytes"
// 	"log"
// 	"os/exec"
// 	"strings"
//
// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/tubbebubbe/transmission"
// )
//
// //	type Torrent struct {
// //		Id       string
// //		Filename string
// //		Status   string
// //	}
//
// type MsgUpdate struct {
// 	Entries []transmission.Torrent
// }
//
// func CmdUpdate() tea.Msg {
//     return
// 	args := []string{"--list"}
// 	cmd := exec.Command("transmission-remote", args...)
// 	stdout, err := cmd.Output()
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	return MsgUpdate{Entries: ParseEntries(stdout)}
// }
//
//
// // func ParseEntries(output []byte) []transmission.Torrent {
// //
// // 	scanner := bufio.NewScanner(bytes.NewReader(output))
// //
// // 	var lines []string
// // 	for scanner.Scan() {
// // 		lines = append(lines, scanner.Text())
// // 	}
// //
// // 	// remove header and "sum" line
// // 	if len(lines) > 2 {
// // 		lines = lines[1 : len(lines)-1]
// // 	}
// //
// // 	var entries []transmission.Torrent
// // 	for _, line := range lines {
// // 		fields := strings.Fields(line)
// // 		_ = log.Default
// // 		fileName := strings.Join(fields[8:], " ")
// // 		_ = fileName
// // 		entries = append(entries, NewEntry(fields[0], fileName, fields[7]))
// // 	}
// //
// // 	return entries
// // }
