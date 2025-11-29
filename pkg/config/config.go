package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds all configuration for zshellcheck.
type Config struct {
	DisabledKatas []string `yaml:"disabled_katas"`

	// Color configuration for text reporter
	ErrorColor   string `yaml:"error_color"`
	WarningColor string `yaml:"warning_color"`
	InfoColor    string `yaml:"info_color"`
	IDColor      string `yaml:"id_color"`
	TitleColor   string `yaml:"title_color"`
	MessageColor string `yaml:"message_color"`
	LineColor    string `yaml:"line_color"`
	ColumnColor  string `yaml:"column_color"`
	NoColor      bool   `yaml:"no_color"`
	Verbose      bool   `yaml:"verbose"`
}

// Default colors
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
	ColorBold   = "\033[1m"
)

const Banner = "\n" +
	"\033[38;5;51m███████╗███████╗██╗  ██╗███████╗██╗     ██╗      ██████╗██╗  ██╗███████╗ ██████╗██╗  ██╗\033[0m\n" +
	"\033[38;5;45m╚══███╔╝██╔════╝██║  ██║██╔════╝██║     ██║     ██╔════╝██║  ██║██╔════╝██╔════╝██║ ██╔╝\033[0m\n" +
	"\033[38;5;39m  ███╔╝ ███████╗███████║█████╗  ██║     ██║     ██║     ███████║█████╗  ██║     █████╔╝\033[0m\n" +
	"\033[38;5;33m ███╔╝  ╚════██║██╔══██║██╔══╝  ██║     ██║     ██║     ██╔══██║██╔══╝  ██║     ██╔═██╗\033[0m\n" +
	"\033[38;5;27m███████╗███████║██║  ██║███████╗███████╗███████╗╚██████╗██║  ██║███████╗╚██████╗██║  ██╗\033[0m\n" +
	"\033[38;5;21m╚══════╝╚══════╝╚═╝  ╚═╝╚══════╝╚══════╝╚══════╝ ╚═════╝╚═╝  ╚═╝╚══════╝ ╚═════╝╚═╝  ╚═╝\033[0m\n" +
	"\n"

// DefaultConfig returns a Config with default values.
func DefaultConfig() Config {
	return Config{
		ErrorColor:   ColorRed,
		WarningColor: ColorYellow,
		InfoColor:    ColorCyan,
		IDColor:      ColorRed,
		TitleColor:   ColorCyan,
		MessageColor: ColorReset,
		LineColor:    ColorCyan,
		ColumnColor:  ColorYellow,
		NoColor:      false,
		Verbose:      false,
	}
}

// MergeConfig merges values from `override` into `base`.
func MergeConfig(base, override Config) Config {
	if len(override.DisabledKatas) > 0 {
		base.DisabledKatas = override.DisabledKatas
	}

	if override.ErrorColor != "" {
		base.ErrorColor = override.ErrorColor
	}
	if override.WarningColor != "" {
		base.WarningColor = override.WarningColor
	}
	if override.InfoColor != "" {
		base.InfoColor = override.InfoColor
	}
	if override.IDColor != "" {
		base.IDColor = override.IDColor
	}
	if override.TitleColor != "" {
		base.TitleColor = override.TitleColor
	}
	if override.MessageColor != "" {
		base.MessageColor = override.MessageColor
	}
	if override.LineColor != "" {
		base.LineColor = override.LineColor
	}
	if override.ColumnColor != "" {
		base.ColumnColor = override.ColumnColor
	}
	// These are boolean flags, direct assignment is fine
	base.NoColor = override.NoColor
	base.Verbose = override.Verbose

	return base
}

// NewConfigFromYAML loads configuration from a YAML file.
func NewConfigFromYAML(path string) (Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	var fileConfig Config
	err = yaml.Unmarshal(data, &fileConfig)
	if err != nil {
		return cfg, err
	}

	return MergeConfig(cfg, fileConfig), nil
}
