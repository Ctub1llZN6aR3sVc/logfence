package transform_test

import (
	"regexp"
	"testing"

	"github.com/yourorg/logfence/internal/transform"
)

func entry(kv ...any) map[string]any {
	m := make(map[string]any, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i].(string)] = kv[i+1]
	}
	return m
}

func TestRedact_KnownField(t *testing.T) {
	r := &transform.RedactTransformer{Fields: []string{"password"}, Mask: "***"}
	e := entry("password", "s3cr3t", "user", "alice")
	r.Apply(e)
	if e["password"] != "***" {
		t.Fatalf("expected masked password, got %v", e["password"])
	}
	if e["user"] != "alice" {
		t.Fatalf("user field should be unchanged")
	}
}

func TestRedact_MissingField(t *testing.T) {
	r := &transform.RedactTransformer{Fields: []string{"token"}, Mask: "***"}
	e := entry("msg", "hello")
	r.Apply(e) // must not panic
	if _, ok := e["token"]; ok {
		t.Fatal("token should not be injected")
	}
}

func TestRename(t *testing.T) {
	rn := &transform.RenameTransformer{Mapping: map[string]string{"msg": "message"}}
	e := entry("msg", "hello")
	rn.Apply(e)
	if e["message"] != "hello" {
		t.Fatalf("expected renamed field, got %v", e["message"])
	}
	if _, ok := e["msg"]; ok {
		t.Fatal("old key should be removed")
	}
}

func TestAddFields(t *testing.T) {
	a := &transform.AddFieldsTransformer{Fields: map[string]any{"env": "prod"}}
	e := entry("msg", "hi")
	a.Apply(e)
	if e["env"] != "prod" {
		t.Fatalf("expected env=prod, got %v", e["env"])
	}
}

func TestRegexRedact(t *testing.T) {
	rx := &transform.RegexRedactTransformer{
		Field:   "msg",
		Pattern: regexp.MustCompile(`\b\d{4}-\d{4}-\d{4}-\d{4}\b`),
		Mask:    "[CARD]",
	}
	e := entry("msg", "charged 1234-5678-9012-3456 ok")
	rx.Apply(e)
	if e["msg"] != "charged [CARD] ok" {
		t.Fatalf("unexpected value: %v", e["msg"])
	}
}

func TestChain(t *testing.T) {
	c := transform.Chain(
		&transform.RedactTransformer{Fields: []string{"secret"}, Mask: "REDACTED"},
		&transform.AddFieldsTransformer{Fields: map[string]any{"host": "box1"}},
	)
	e := entry("secret", "abc", "msg", "test")
	c.Apply(e)
	if e["secret"] != "REDACTED" {
		t.Fatalf("chain: redact failed")
	}
	if e["host"] != "box1" {
		t.Fatalf("chain: add fields failed")
	}
}

func TestNormalizeLevel(t *testing.T) {
	e := entry("level", "ERROR")
	transform.NormalizeLevel(e)
	if e["level"] != "error" {
		t.Fatalf("expected lowercase level, got %v", e["level"])
	}
}
