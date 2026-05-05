package filter

import (
	"testing"
)

func TestParseLevel(t *testing.T) {
	cases := []struct {
		input string
		want  Level
		ok    bool
	}{
		{"debug", LevelDebug, true},
		{"INFO", LevelInfo, true},
		{"Warn", LevelWarn, true},
		{"ERROR", LevelError, true},
		{"trace", 0, false},
	}
	for _, tc := range cases {
		got, ok := ParseLevel(tc.input)
		if ok != tc.ok {
			t.Errorf("ParseLevel(%q) ok=%v, want %v", tc.input, ok, tc.ok)
		}
		if ok && got != tc.want {
			t.Errorf("ParseLevel(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestRuleMatch_Level(t *testing.T) {
	r := Rule{MinLevel: LevelWarn}

	if r.Match(Entry{Level: LevelDebug}) {
		t.Error("debug entry should not match warn rule")
	}
	if r.Match(Entry{Level: LevelInfo}) {
		t.Error("info entry should not match warn rule")
	}
	if !r.Match(Entry{Level: LevelWarn}) {
		t.Error("warn entry should match warn rule")
	}
	if !r.Match(Entry{Level: LevelError}) {
		t.Error("error entry should match warn rule")
	}
}

func TestRuleMatch_Fields(t *testing.T) {
	r := Rule{
		MinLevel: LevelDebug,
		Fields:   map[string]string{"service": "api"},
	}

	matching := Entry{Level: LevelInfo, Fields: map[string]string{"service": "api"}}
	if !r.Match(matching) {
		t.Error("entry with matching field should pass")
	}

	nonMatching := Entry{Level: LevelInfo, Fields: map[string]string{"service": "worker"}}
	if r.Match(nonMatching) {
		t.Error("entry with non-matching field should not pass")
	}
}

func TestChain_Empty(t *testing.T) {
	e := Entry{Level: LevelDebug}
	if !Chain(nil, e) {
		t.Error("empty chain should match everything")
	}
}

func TestChain_OR(t *testing.T) {
	rules := []Rule{
		{MinLevel: LevelError},
		{MinLevel: LevelDebug, Fields: map[string]string{"service": "api"}},
	}

	// Matches second rule
	e1 := Entry{Level: LevelInfo, Fields: map[string]string{"service": "api"}}
	if !Chain(rules, e1) {
		t.Error("e1 should match via second rule")
	}

	// Matches first rule
	e2 := Entry{Level: LevelError, Fields: map[string]string{"service": "db"}}
	if !Chain(rules, e2) {
		t.Error("e2 should match via first rule")
	}

	// Matches neither
	e3 := Entry{Level: LevelInfo, Fields: map[string]string{"service": "db"}}
	if Chain(rules, e3) {
		t.Error("e3 should not match any rule")
	}
}
