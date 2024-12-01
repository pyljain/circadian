package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Targets []Target
}

type Target struct {
	Name            string            `yaml:"name"`
	URL             string            `yaml:"url"`
	Method          string            `yaml:"method"`
	Headers         map[string]string `yaml:"headers"`
	Body            string            `yaml:"body"`
	Interval        int               `yaml:"interval"`
	ExpectedLatency int               `yaml:"expected_latency"`
}

func Read(path string) (*Config, error) {
	configBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := Config{}
	err = yaml.Unmarshal(configBytes, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil

}
