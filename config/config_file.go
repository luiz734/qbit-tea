package config

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/go-playground/validator/v10"
	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	// Dirs user can download
	MoviesDirs []string `toml:"movies_dirs" validate:"required"`
	ShowsDirs  []string `toml:"shows_dirs" validate:"required"`
}

func ReadConfigFile(cli CLI) (*Config, error) {
	var cfg Config

	// Try to read the config file
	file, err := os.ReadFile(cli.ConfigFile)
	if err != nil {
		return nil, fmt.Errorf("can't read file %s: %w", cli.ConfigFile, err)
	}
	log.Debug("Read default config file", "path", cli.ConfigFile)

	// Config file should be avaliabe now
	err = toml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, fmt.Errorf("can't parse config file: %w: ", err)
	}
	log.Debug("Config file sucessfully unmarshaled")

	// Validate structs
	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("can't validate config file: %w", err)
	}
	log.Debug("Config file sucessfully validated")

	log.Info("Config file ok")
	return &cfg, nil

}

func validateConfig(cfg Config) error {
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return fmt.Errorf("validation failed: %w: ", err)
	}
	// err := validateWidgetsConfig(cfg.Widgets)
	return nil
}
