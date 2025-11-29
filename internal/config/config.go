package config

type PathConfig struct {
	Source string `yaml:"source"`
	Type   string `yaml:"type"` // "tile" or "directory"
	Name   string `yaml:"name"`
}

type Config struct {
	DotfilesDir string       `yaml:"dotfiles_dir"`
	GitRepo     string       `yaml:"git_repo"`
	Paths       []PathConfig `yaml:"paths"`
}
