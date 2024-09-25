package main

import (
	"fmt"
	"log"
	"os"
	"qbit-tea/app"
	"qbit-tea/input"
	"time"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/tubbebubbe/transmission"
)

type actionMsg struct {
	helpItem input.UserAction
}

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}
	client := transmission.New("http://127.0.0.1:9091", "", "")
	// util.CheckError(err)
	_, err := client.GetTorrents()
	if err != nil {
		log.Panic(err)
	}

	program := tea.NewProgram(app.NewModel(timer.NewWithInterval(app.Timeout, time.Millisecond), &client))
	if err != nil {

		fmt.Printf("Uh oh, there was an error: %v\n", err)
	}
	program.Run()

}
