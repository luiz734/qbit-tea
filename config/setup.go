package config

import (
	"fmt"
	"github.com/charmbracelet/log"
	"os"
	"path/filepath"
)

func CreateConfigFile(cli *CLI) error {

	// Read user
	user := os.Getenv("USER")
	if user == "" {
		return fmt.Errorf("can't read env $USER")
	}
	log.Debug("Read environment variable", "$USER", user)

	// Set the config file path
	defaultConfigPath := filepath.Join("/home", user, ".config", "qbit-tea", "config.toml")
	// User provide a config file? use it
	if cli.ConfigFile == "" {
		cli.ConfigFile = defaultConfigPath
	}

	configDir := filepath.Dir(cli.ConfigFile)

	if err := createDirIfNotExists(configDir); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	log.Debug("Config dir ok", "dir", configDir)

	// File not exists
	if _, err := os.Stat(cli.ConfigFile); os.IsNotExist(err) {
		if err := os.WriteFile(cli.ConfigFile, defaultConfigToml, 0755); err != nil {
			return fmt.Errorf("can't create file at %s: %w", cli.ConfigFile, err)
		}
	} else {
		log.Info("Found config file")
	}

	return nil
}

func createDirIfNotExists(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return fmt.Errorf("can't create config dir at %s: %w", dirPath, err)
		}
	}
	log.Info("Create dir", "dirpath", dirPath)
	return nil
}

func createFileIfNotExists(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := os.WriteFile(filePath, []byte(""), 0755); err != nil {
			return fmt.Errorf("can't create file at %s: %w", filePath, err)
		}
	}
	log.Info("Create file", "filepath", filePath)
	return nil
}
