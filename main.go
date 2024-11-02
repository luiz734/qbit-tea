package main

import (
	"fmt"
	"io"
	"os"
	"qbit-tea/app"
	"qbit-tea/config"
	"time"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"

	"github.com/tubbebubbe/transmission"
)

func main() {
	var err error

	// Setup logs
	logOutputFile := setupLogs()
	defer logOutputFile.Close()

	var cli config.CLI
	_ = kong.Parse(&cli)
	address := fmt.Sprintf("http://%s", cli.Address)

	// Create the config file
	// Update the cli.ConfigFile to the default if no was provided
	if err = config.CreateConfigFile(&cli); err != nil {
		log.Fatalf("Error creating config file: %v", err)
	}
	if config.Cfg, err = config.ReadConfigFile(cli); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	log.Info("Found", "movies", config.Cfg.MoviesDirs)
	log.Info("Found", "shows", config.Cfg.ShowsDirs)

	log.Info("Trying connect ", "addres", address)
	client := transmission.New(address, cli.User, cli.Password)
	// util.CheckError(err)
	_, err = client.GetTorrents()
	if err != nil {
		log.Fatal("Can't connect to transmission-daemon.\nIs the daemon running?\n%v", err)
	}

	program := tea.NewProgram(app.NewModel(timer.NewWithInterval(app.Timeout, time.Millisecond), &client, cli.Address))
	// if err != nil {
	// 	log.Fatal("Uh oh, there was an error: %v\n", err)
	// }
	program.Run()

}
func setupLogs() *os.File {
	var f *os.File

	// Redirect log output to file
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
		log.SetOutput(f)
	} else {
		log.SetOutput(io.Discard)
	}

	// if len(os.Getenv("DEBUG")) > 0 {
	//
	// 	logger := log.NewWithOptions(os.Stderr, log.Options{
	// 		ReportTimestamp: false,
	// 		ReportCaller:    true,
	// 		Level:           log.DebugLevel,
	// 	})
	// 	log.SetDefault(logger)
	//
	// 	f, err = os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// 	if err != nil {
	// 		log.Fatalf("Error opening file: %v", err)
	// 	}
	// 	// Timestamp only on file
	// 	log.SetReportTimestamp(true)
	// 	log.SetReportCaller(false)
	// 	log.SetOutput(f)
	// 	log.Printf("Set log output to file %s", "debug.log")
	// } else {
	// 	// _ = io.Discard
	// 	log.SetOutput(io.Discard)
	// }
	return f
}
