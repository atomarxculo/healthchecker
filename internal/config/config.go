package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Service struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Target string `yaml:"target"`
	Method string `yaml:"method"`
	Timeout time.Duration `yaml:"timeout"`
}

type Config struct {
	Services []Service `yaml:"services"`
	Global struct {
		Timeout time.Duration `yaml:"timeout"`
	} `yaml:"global"`
}

func Load(path string) (*Config, error)  {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil{
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	for i, svc := range cfg.Services{
		if svc.Name == "" {
			return nil, fmt.Errorf("service at index %d has no name", i)
		}
		if svc.Type == "" {
			return nil, fmt.Errorf("service %s has no type", svc.Name)
		}
		if svc.Target == "" {
			return nil, fmt.Errorf("service %s has no target", svc.Name)
		}
	}

	return &cfg, nil

}