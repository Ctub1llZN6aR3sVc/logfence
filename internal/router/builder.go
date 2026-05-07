package router

import (
	"fmt"

	"github.com/yourorg/logfence/internal/config"
	"github.com/yourorg/logfence/internal/filter"
	"github.com/yourorg/logfence/internal/ratelimit"
	"github.com/yourorg/logfence/internal/sink"
)

// Build constructs a Router from the loaded configuration.
func Build(cfg *config.Config) (*Router, error) {
	sinks := make(map[string]sink.Sink, len(cfg.Sinks))
	for _, sc := range cfg.Sinks {
		s, err := sink.New(sc.Type, sc.Path)
		if err != nil {
			return nil, fmt.Errorf("builder: sink %q: %w", sc.Name, err)
		}
		sinks[sc.Name] = s
	}

	routes := make([]Route, 0, len(cfg.Routes))
	for _, rc := range cfg.Routes {
		chain, err := buildChain(rc.Filters)
		if err != nil {
			return nil, fmt.Errorf("builder: route %q: %w", rc.Name, err)
		}

		s, ok := sinks[rc.Sink]
		if !ok {
			return nil, fmt.Errorf("builder: route %q: sink %q not found", rc.Name, rc.Sink)
		}

		var lim *ratelimit.Limiter
		if rc.RateLimit.Rate > 0 || rc.RateLimit.Burst > 0 {
			burst := rc.RateLimit.Burst
			if burst == 0 {
				burst = rc.RateLimit.Rate
			}
			lim = ratelimit.New(rc.RateLimit.Rate, burst)
		}

		routes = append(routes, Route{
			Name:    rc.Name,
			Chain:   chain,
			Sink:    s,
			Limiter: lim,
		})
	}

	return New(routes), nil
}

func buildChain(specs []string) (filter.Chain, error) {
	rules := make([]filter.Rule, 0, len(specs))
	for _, spec := range specs {
		r, err := buildRule(spec)
		if err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	return filter.Chain(rules), nil
}

func buildRule(spec string) (filter.Rule, error) {
	return filter.ParseRule(spec)
}
