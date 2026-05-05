package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Route defines a log routing rule: match logs by level/source and forward to a sink.
type Route struct {
	Name    string   `yaml:"name"`
	Levels  []string `yaml:"levels"`
	Sources []string `yaml:"sources"`
	Sink    string   `yaml:"sink"`
}

// Sink defines an output destination for log entries.
type Sink struct {
	Name   string            `yaml:"name"`
	Type   string            `yaml:"type"` // stdout, file, http
	Target string            `yaml:"target"`
	Opts   map[string]string `yaml:"opts"`
}

// Config is the top-level logfence configuration.
type Config struct {
	ListenAddr string  `yaml:"listen_addr"`
	Routes     []Route `yaml:"routes"`
	Sinks      []Sink  `yaml:"sinks"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse yaml: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config: validation: %w", err)
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.ListenAddr == "" {
		c.ListenAddr = ":5170"
	}

	sinkNames := make(map[string]struct{}, len(c.Sinks))
	for _, s := range c.Sinks {
		if s.Name == "" {
			return fmt.Errorf("sink missing name")
		}
		if s.Type == "" {
			return fmt.Errorf("sink %q missing type", s.Name)
		}
		sinkNames[s.Name] = struct{}{}
	}

	for _, r := range c.Routes {
		if r.Sink == "" {
			return fmt.Errorf("route %q missing sink", r.Name)
		}
		if _, ok := sinkNames[r.Sink]; !ok {
			return fmt.Errorf("route %q references unknown sink %q", r.Name, r.Sink)
		}
	}

	return nil
}
