package config

import (
	"errors"
	"fmt"
)

func (c *Config) validate() error {
	if err := c.Server.validate(); err != nil {
		return err
	}
	for name, provider := range c.LLMProviders {
		if err := provider.validate(); err != nil {
			return fmt.Errorf("provider '%s' config error: %w", name, err)
		}
	}
	return nil
}

func (sc ServerConfig) validate() error {
	if sc.Port <= 0 || sc.Port > 65535 {
		return fmt.Errorf("server port must be between 1 and 65535, but got %d", sc.Port)
	}
	if sc.ReadTimeout <= 0 {
		return errors.New("server read_timeout must be strictly positive")
	}
	if sc.WriteTimeout <= 0 {
		return errors.New("server write_timeout must be strictly positive")
	}
	if sc.IdleTimeout <= 0 {
		return errors.New("server idle_timeout must be strictly positive")
	}
	if sc.RequestContextTimeout <= 0 {
		return errors.New("server request_context_timeout must be strictly positive")
	}
	if sc.MaxHeaderBytes <= 0 {
		return errors.New("server max_header_bytes must be strictly positive")
	}
	return nil
}

func (lpc LLMProviderConfig) validate() error {
	if !lpc.Enabled {
		return nil
	}
	if lpc.APIKey == "" {
		return errors.New("api_key is required")
	}
	if lpc.Model == "" {
		return errors.New("model name is required")
	}
	return nil
}
