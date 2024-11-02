package config

type CLI struct {
	Address    string `short:"a" name:"address" default:"localhost:9091" help:"Address"`
	User       string `short:"u" name:"user" default:"" help:"Transmission user"`
	Password   string `short:"p" name:"password" default:"" help:"Transmission password"`
	ConfigFile string `short:"c" name:"config" default:"" help:"Config file path\nDefaults to ~/.config/qbit-tea/config.toml"`
}

var Cfg *Config
