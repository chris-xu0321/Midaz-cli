package output

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestWriteSuccess_EmptyMeta(t *testing.T) {
	var buf bytes.Buffer
	err := WriteSuccess(&buf, map[string]string{"status": "ok"}, nil, "json")
	if err != nil {
		t.Fatal(err)
	}
	var result map[string]any
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result["ok"] != true {
		t.Error("expected ok=true")
	}
	meta, ok := result["meta"].(map[string]any)
	if !ok {
		t.Fatal("meta should be an object, got nil or wrong type")
	}
	if len(meta) != 0 {
		t.Errorf("expected empty meta, got %v", meta)
	}
}

func TestWriteSuccess_ZeroValueIntPreserved(t *testing.T) {
	var buf bytes.Buffer
	meta := map[string]any{
		"count":              5,
		"contradicting_count": 0,
	}
	err := WriteSuccess(&buf, []string{"a"}, meta, "json")
	if err != nil {
		t.Fatal(err)
	}
	var result map[string]any
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	m := result["meta"].(map[string]any)
	if m["contradicting_count"] != float64(0) {
		t.Errorf("expected contradicting_count=0, got %v", m["contradicting_count"])
	}
	if m["count"] != float64(5) {
		t.Errorf("expected count=5, got %v", m["count"])
	}
}

func TestWriteSuccess_PrettyFormat(t *testing.T) {
	var buf bytes.Buffer
	err := WriteSuccess(&buf, "hello", nil, "pretty")
	if err != nil {
		t.Fatal(err)
	}
	// Pretty output should contain newlines and indentation
	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("\n  ")) {
		t.Error("expected indented output for pretty format")
	}
}

func TestWriteRaw_AddsNewline(t *testing.T) {
	var buf bytes.Buffer
	err := WriteRaw(&buf, []byte(`{"status":"ok"}`), "json")
	if err != nil {
		t.Fatal(err)
	}
	if buf.Bytes()[buf.Len()-1] != '\n' {
		t.Error("expected trailing newline")
	}
}
