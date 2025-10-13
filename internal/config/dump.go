package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func Dump() error {
	out, err := yaml.Marshal(Default)
	if err != nil {
		return fmt.Errorf("marshal config error: %w", err)
	}
	fmt.Printf("%s", out)
	return nil
}
