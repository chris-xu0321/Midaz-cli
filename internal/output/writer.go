package output

import (
	"bytes"
	"encoding/json"
	"io"
)

// WriteSuccess writes a success envelope to w.
// Meta uses map[string]any so zero-value ints (e.g., "contradicting_count": 0) are preserved.
func WriteSuccess(w io.Writer, data interface{}, meta map[string]any, format string) error {
	if meta == nil {
		meta = map[string]any{}
	}
	env := Envelope{
		OK:   true,
		Data: data,
		Meta: meta,
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	if format == "pretty" {
		enc.SetIndent("", "  ")
	}
	return enc.Encode(env)
}

// WriteRaw writes raw API response bytes to w.
// If format is "pretty", re-indents the JSON.
func WriteRaw(w io.Writer, raw []byte, format string) error {
	if format == "pretty" {
		var buf bytes.Buffer
		if err := json.Indent(&buf, raw, "", "  "); err != nil {
			// If not valid JSON, write as-is
			_, err := w.Write(raw)
			return err
		}
		buf.WriteByte('\n')
		_, err := buf.WriteTo(w)
		return err
	}
	// Compact: write as-is with trailing newline
	if _, err := w.Write(raw); err != nil {
		return err
	}
	if len(raw) > 0 && raw[len(raw)-1] != '\n' {
		_, err := w.Write([]byte{'\n'})
		return err
	}
	return nil
}
