package config

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

func Load(configPath string) (Config, error) {
	k := koanf.New(".")
	var cfg Config
	if err := loadConfigFile(k, configPath); err != nil {
		return cfg, err
	}
	if err := k.Unmarshal("", &cfg); err != nil {
		return cfg, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	if err := cfg.validate(); err != nil {
		return cfg, fmt.Errorf("invalid configuration: %w", err)
	}
	cfg.Logger = newLogger(cfg.Log)
	return cfg, nil
}

func loadConfigFile(k *koanf.Koanf, configPath string) error {
	err := k.Load(file.Provider(configPath), yaml.Parser())
	if err != nil {
		return fmt.Errorf("failed to load config file '%s': %w", configPath, err)
	}
	return nil
}

func newLogger(cfg LogConfig) zerolog.Logger {
	level, err := zerolog.ParseLevel(strings.ToLower(cfg.Level))
	if err != nil {
		level = zerolog.InfoLevel
	}
	var output io.Writer = os.Stderr
	if cfg.Pretty {
		output = zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		}
	}
	return zerolog.New(output).Level(level).With().Timestamp().Logger()
}
