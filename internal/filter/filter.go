package filter

import (
	"strings"
)

// Level represents a log severity level.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// levelNames maps string representations to Level values.
var levelNames = map[string]Level{
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"error": LevelError,
}

// ParseLevel converts a string to a Level, returning an error if unknown.
func ParseLevel(s string) (Level, bool) {
	l, ok := levelNames[strings.ToLower(s)]
	return l, ok
}

// Rule defines a single filter rule with a minimum log level and optional
// field-based inclusion constraints.
type Rule struct {
	MinLevel Level
	// Fields is an optional map of field key→value that must all match.
	Fields map[string]string
}

// Entry represents a parsed log entry passed through the filter.
type Entry struct {
	Level  Level
	Fields map[string]string
	Raw    string
}

// Match reports whether the entry satisfies the rule.
func (r *Rule) Match(e Entry) bool {
	if e.Level < r.MinLevel {
		return false
	}
	for k, v := range r.Fields {
		if e.Fields[k] != v {
			return false
		}
	}
	return true
}

// Chain applies a slice of rules to an entry, returning true if ANY rule
// matches (logical OR). An empty chain matches everything.
func Chain(rules []Rule, e Entry) bool {
	if len(rules) == 0 {
		return true
	}
	for _, r := range rules {
		if r.Match(e) {
			return true
		}
	}
	return false
}
