// Package transform provides log entry field transformation
// capabilities such as redacting sensitive fields, renaming keys,
// and adding static metadata before routing to sinks.
package transform

import (
	"regexp"
	"strings"
)

// Transformer mutates a log entry map in-place.
type Transformer interface {
	Apply(entry map[string]any)
}

// RedactTransformer replaces the value of matched fields with a
// fixed mask string.
type RedactTransformer struct {
	Fields []string
	Mask   string
}

func (r *RedactTransformer) Apply(entry map[string]any) {
	for _, f := range r.Fields {
		if _, ok := entry[f]; ok {
			entry[f] = r.Mask
		}
	}
}

// RenameTransformer renames entry keys according to the provided map.
type RenameTransformer struct {
	// From -> To
	Mapping map[string]string
}

func (rn *RenameTransformer) Apply(entry map[string]any) {
	for from, to := range rn.Mapping {
		if v, ok := entry[from]; ok {
			entry[to] = v
			delete(entry, from)
		}
	}
}

// AddFieldsTransformer injects static key/value pairs into every entry.
type AddFieldsTransformer struct {
	Fields map[string]any
}

func (a *AddFieldsTransformer) Apply(entry map[string]any) {
	for k, v := range a.Fields {
		entry[k] = v
	}
}

// RegexRedactTransformer redacts substrings in string field values
// that match a regular expression.
type RegexRedactTransformer struct {
	Field   string
	Pattern *regexp.Regexp
	Mask    string
}

func (rx *RegexRedactTransformer) Apply(entry map[string]any) {
	v, ok := entry[rx.Field]
	if !ok {
		return
	}
	s, ok := v.(string)
	if !ok {
		return
	}
	entry[rx.Field] = rx.Pattern.ReplaceAllString(s, rx.Mask)
}

// Chain applies a sequence of Transformers in order.
func Chain(transformers ...Transformer) Transformer {
	return &chainTransformer{steps: transformers}
}

type chainTransformer struct {
	steps []Transformer
}

func (c *chainTransformer) Apply(entry map[string]any) {
	for _, t := range c.steps {
		t.Apply(entry)
	}
}

// NormalizeLevel lower-cases the "level" field if present.
func NormalizeLevel(entry map[string]any) {
	if v, ok := entry["level"]; ok {
		if s, ok := v.(string); ok {
			entry["level"] = strings.ToLower(s)
		}
	}
}
