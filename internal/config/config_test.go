package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestConfig_YAMLUnmarshal(t *testing.T) {
	raw := `
dotfiles_dir: dotfiles
git_repo: https://github.com/example/dotfiles.git
paths:
  - source: ~/.zshrc
    type: file
    name: zsh
  - source: ~/.config/nvim
    type: directory
    name: nvim
`
	var cfg Config
	err := yaml.Unmarshal([]byte(raw), &cfg)
	require.NoError(t, err)

	assert.Equal(t, "dotfiles", cfg.DotfilesDir)
	assert.Equal(t, "https://github.com/example/dotfiles.git", cfg.GitRepo)
	require.Len(t, cfg.Paths, 2)
	assert.Equal(t, "~/.zshrc", cfg.Paths[0].Source)
	assert.Equal(t, "file", cfg.Paths[0].Type)
	assert.Equal(t, "zsh", cfg.Paths[0].Name)
	assert.Equal(t, "~/.config/nvim", cfg.Paths[1].Source)
	assert.Equal(t, "directory", cfg.Paths[1].Type)
	assert.Equal(t, "nvim", cfg.Paths[1].Name)
}

func TestPathConfig_Empty(t *testing.T) {
	raw := `source: ""
type: file
name: ""
`
	var p PathConfig
	err := yaml.Unmarshal([]byte(raw), &p)
	require.NoError(t, err)
	assert.Empty(t, p.Source)
	assert.Empty(t, p.Name)
}
