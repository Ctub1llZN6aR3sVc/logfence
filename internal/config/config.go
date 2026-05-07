// Package config loads and validates logfence configuration from a YAML file.
package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// RateLimit holds per-route rate limiting parameters.
type RateLimit struct {
	Rate  int `yaml:"rate"`  // tokens per second
	Burst int `yaml:"burst"` // maximum burst size
}

// Route maps a filter chain to a sink, with optional rate limiting.
type Route struct {
	Name      string    `yaml:"name"`
	Filters   []string  `yaml:"filters"`
	Sink      string    `yaml:"sink"`
	RateLimit RateLimit `yaml:"rate_limit"`
}

// Sink defines an output destination.
type Sink struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Path string `yaml:"path,omitempty"`
}

// Config is the top-level logfence configuration.
type Config struct {
	ListenAddr string  `yaml:"listen_addr"`
	Routes     []Route `yaml:"routes"`
	Sinks      []Sink  `yaml:"sinks"`
}

// Load reads, parses, and validates the configuration at path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse: %w", err)
	}

	if cfg.ListenAddr == "" {
		cfg.ListenAddr = ":8080"
	}

	sinkNames := make(map[string]struct{}, len(cfg.Sinks))
	for _, s := range cfg.Sinks {
		if s.Type != "stdout" && s.Type != "file" {
			return nil, fmt.Errorf("config: unknown sink type %q", s.Type)
		}
		sinkNames[s.Name] = struct{}{}
	}

	for _, r := range cfg.Routes {
		if _, ok := sinkNames[r.Sink]; !ok {
			return nil, fmt.Errorf("config: route %q references unknown sink %q", r.Name, r.Sink)
		}
		if err := validateRateLimit(r.RateLimit); err != nil {
			return nil, fmt.Errorf("config: route %q rate_limit: %w", r.Name, err)
		}
	}

	return &cfg, nil
}

func validateRateLimit(rl RateLimit) error {
	if rl.Rate < 0 {
		return errors.New("rate must be >= 0")
	}
	if rl.Burst < 0 {
		return errors.New("burst must be >= 0")
	}
	return nil
}
