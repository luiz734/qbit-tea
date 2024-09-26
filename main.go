package main

import (
	"fmt"
	"log"
	"os"
	"qbit-tea/app"
	"qbit-tea/input"
	"time"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/tubbebubbe/transmission"
)

type actionMsg struct {
	helpItem input.UserAction
}

type CLI struct {
	Address  string `short:"a" name:"address" default:"localhost:9091" help:"Address"`
	User     string `short:"u" name:"user" default:"" help:"Transmission user"`
	Password string `short:"p" name:"password" default:"" help:"Transmission password"`
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

	var cli CLI
	_ = kong.Parse(&cli)
	address := fmt.Sprintf("http://%s", cli.Address)

	client := transmission.New(address, cli.User, cli.Password)
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
