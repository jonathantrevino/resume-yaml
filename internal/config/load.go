package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

// unmarshal config.yaml
func (c *Config) Load(p string) error {

	data, err := os.ReadFile(p)
	if err != nil {
		return fmt.Errorf("reading config file %s: %w", p, err)
	}

	if err := yaml.Unmarshal(data, c); err != nil {
		return fmt.Errorf("unmarshalling config: %w", err)
	}

	return err
}
