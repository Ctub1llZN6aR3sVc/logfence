package sink

import (
	"encoding/json"
	"fmt"
	"io"
)

// writeJSON serialises entry as a single JSON line to w.
func writeJSON(w io.Writer, entry map[string]any) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("sink marshal: %w", err)
	}
	data = append(data, '\n')
	_, err = w.Write(data)
	return err
}
