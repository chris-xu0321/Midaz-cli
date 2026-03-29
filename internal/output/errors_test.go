package output

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestExitError_Error(t *testing.T) {
	err := ErrValidation("missing arg: %s", "query")
	if err.Error() != "missing arg: query" {
		t.Errorf("expected 'missing arg: query', got %q", err.Error())
	}
	if err.Code != ExitValidation {
		t.Errorf("expected exit code %d, got %d", ExitValidation, err.Code)
	}
}

func TestWriteErrorEnvelope_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	err := &ExitError{
		Code:   ExitAPI,
		Detail: &ErrDetail{Code: "not_found", Message: "Thread not found", Hint: "try: seer-q threads"},
	}
	WriteErrorEnvelope(&buf, err)

	var result map[string]any
	if jsonErr := json.Unmarshal(buf.Bytes(), &result); jsonErr != nil {
		t.Fatalf("invalid JSON: %v\nOutput: %s", jsonErr, buf.String())
	}
	if result["ok"] != false {
		t.Error("expected ok=false")
	}
	errObj := result["error"].(map[string]any)
	if errObj["code"] != "not_found" {
		t.Errorf("expected code=not_found, got %v", errObj["code"])
	}
	if errObj["hint"] != "try: seer-q threads" {
		t.Errorf("expected hint, got %v", errObj["hint"])
	}
}

func TestErrNetwork(t *testing.T) {
	err := ErrNetwork("connection refused to %s", "localhost")
	if err.Code != ExitNetwork {
		t.Errorf("expected exit code %d, got %d", ExitNetwork, err.Code)
	}
	if err.Detail.Code != "network" {
		t.Errorf("expected error code 'network', got %q", err.Detail.Code)
	}
}
