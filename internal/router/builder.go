package router

import (
	"fmt"

	"github.com/yourorg/logfence/internal/config"
	"github.com/yourorg/logfence/internal/filter"
	"github.com/yourorg/logfence/internal/sink"
)

// Build constructs a Router from the application Config.
// It wires together filter chains and sinks for every configured route.
func Build(cfg *config.Config) (*Router, error) {
	var routes []Route

	for _, rc := range cfg.Routes {
		// Resolve sink.
		sinkCfg, ok := cfg.Sinks[rc.Sink]
		if !ok {
			return nil, fmt.Errorf("builder: route %q references unknown sink %q", rc.Name, rc.Sink)
		}

		s, err := sink.New(sinkCfg)
		if err != nil {
			return nil, fmt.Errorf("builder: route %q sink error: %w", rc.Name, err)
		}

		// Build filter chain.
		var rules []filter.Rule
		for _, rf := range rc.Filters {
			rule, err := buildRule(rf)
			if err != nil {
				return nil, fmt.Errorf("builder: route %q filter error: %w", rc.Name, err)
			}
			rules = append(rules, rule)
		}

		routes = append(routes, Route{
			Name:   rc.Name,
			Filter: filter.Chain{Rules: rules},
			Sink:   s,
		})
	}

	return New(routes), nil
}

func buildRule(rf config.FilterConfig) (filter.Rule, error) {
	var rule filter.Rule

	if rf.Level != "" {
		lvl, err := filter.ParseLevel(rf.Level)
		if err != nil {
			return rule, err
		}
		rule.MinLevel = lvl
	}

	if len(rf.Fields) > 0 {
		rule.Fields = rf.Fields
	}

	return rule, nil
}
