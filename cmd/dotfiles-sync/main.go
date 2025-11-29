package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/mordvn/dotfiles-sync/internal/config"
	"github.com/mordvn/dotfiles-sync/internal/ui"
	"gopkg.in/yaml.v3"
)

func main() {
	configPath := "config.yaml"
	if _, err := os.Stat(configPath); err != nil {
		log.Fatal("config file not found", "err", err)
	}

	cfgData, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal("failed to read config file", "err", err)
	}

	var cfg config.Config
	if err := yaml.Unmarshal(cfgData, &cfg); err != nil {
		log.Fatal("failed to parse config file", "err", err)
	}

	m := ui.NewModel(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal("failed to run program", "err", err)
	}
}
