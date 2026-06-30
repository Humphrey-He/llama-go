package config

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server  ServerConfig            `yaml:"server"`
	Plugins map[string]PluginConfig `yaml:"plugins"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type PluginConfig struct {
	Enabled bool                   `yaml:"enabled"`
	BaseURL string                 `yaml:"base_url"`
	Models  []string               `yaml:"models"`
	Extra   map[string]interface{} `yaml:",inline"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, nil
}
