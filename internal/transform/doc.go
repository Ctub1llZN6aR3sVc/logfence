// Package transform provides composable log-entry transformers for
// logfence pipelines.
//
// Transformers operate on a parsed log entry represented as
// map[string]any and mutate it in-place before the entry is forwarded
// to a sink.
//
// Available transformers:
//
//   - RedactTransformer   – replaces sensitive field values with a mask.
//   - RenameTransformer   – renames entry keys.
//   - AddFieldsTransformer – injects static key/value pairs.
//   - RegexRedactTransformer – masks regex-matched substrings in a field.
//
// Transformers can be composed with Chain:
//
//	c := transform.Chain(
//		&transform.RedactTransformer{Fields: []string{"password"}, Mask: "***"},
//		&transform.AddFieldsTransformer{Fields: map[string]any{"env": "prod"}},
//	)
//	c.Apply(entry)
package transform
