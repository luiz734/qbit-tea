package main

import (
	"fmt"
	"io"
	"os"
	"qbit-tea/app/models/errorscreen"
	"qbit-tea/app/models/torrents"
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

	log.Info("Found config", "download_dirs", config.Cfg.DownloadDirs)

	log.Info("Trying connect ", "addres", address)
	client := transmission.New(address, cli.User, cli.Password)
	// util.CheckError(err)
	_, err = client.GetTorrents()

	var program *tea.Program
	if err != nil {
		log.Errorf("Can't connect to transmission-daemon.\nIs the daemon running?\n%v", err)
		errMsg := "Error launching program"
		errDesc := "Is the daemon running?\nIs the address correct?"
		program = tea.NewProgram(errorscreen.InitialModel(errorscreen.QuitModel(), errMsg, fmt.Errorf(errDesc), 0, 0))
		// program = tea.NewProgram(addtorrent.InitialModel(nil))
	} else {
		program = tea.NewProgram(torrents.NewModel(timer.NewWithInterval(torrents.Timeout, time.Millisecond), &client, cli.Address))
	}
	program.Run()

}

// Return the file logs will be redirected to
// Caller can call defer f.Close()
func setupLogs() *os.File {
	var f *os.File

	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		ReportCaller:    true,
		Level:           log.DebugLevel,
	})
	log.SetDefault(logger)

	// Redirect log output to file
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
		}
		log.SetOutput(f)
	} else {
		log.SetOutput(io.Discard)
	}

	return f
}
