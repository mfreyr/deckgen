package config

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type Config struct {
	Server       ServerConfig                 `koanf:"server" yaml:"server"`
	Log          LogConfig                    `koanf:"log" yaml:"log"`
	LLMProviders map[string]LLMProviderConfig `koanf:"llm_providers" yaml:"llm_providers"`
	Logger       zerolog.Logger               `koanf:"-" yaml:"-"`
}
type LLMProviderConfig struct {
	Enabled bool   `koanf:"enabled" yaml:"enabled"`
	APIKey  string `koanf:"api_key" yaml:"api_key"`
	Model   string `koanf:"model" yaml:"model"`
}

type ServerConfig struct {
	Port                  int           `koanf:"port" yaml:"port"`
	ReadTimeout           time.Duration `koanf:"read_timeout" yaml:"read_timeout"`
	WriteTimeout          time.Duration `koanf:"write_timeout" yaml:"write_timeout"`
	IdleTimeout           time.Duration `koanf:"idle_timeout" yaml:"idle_timeout"`
	RequestContextTimeout time.Duration `koanf:"request_context_timeout" yaml:"request_context_timeout"`
	MaxHeaderBytes        int           `koanf:"max_header_bytes" yaml:"max_header_bytes"`
}

type LogConfig struct {
	Level  string `koanf:"level" yaml:"level"`
	Pretty bool   `koanf:"pretty" yaml:"pretty"`
}

var Default = Config{
	Server: ServerConfig{
		Port:                  8080,
		ReadTimeout:           6 * time.Second,
		WriteTimeout:          6 * time.Second,
		IdleTimeout:           120 * time.Second,
		RequestContextTimeout: 8 * time.Second,
		MaxHeaderBytes:        http.DefaultMaxHeaderBytes,
	},
	Log: LogConfig{
		Level:  "info",
		Pretty: false,
	},
	LLMProviders: map[string]LLMProviderConfig{
		"openai": {
			Model: "gpt-5-mini",
		},
	},
}
