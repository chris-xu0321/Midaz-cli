package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/chris-xu0321/Midaz-cli/internal/output"
)

func TestHandleRootError_ExitError(t *testing.T) {
	var buf bytes.Buffer
	err := output.ErrAPI("not_found", "Thread not found")
	code := handleRootError(&buf, err)
	if code != output.ExitAPI {
		t.Errorf("expected exit code %d, got %d", output.ExitAPI, code)
	}
	// Verify JSON envelope on stderr
	var result map[string]any
	if jsonErr := json.Unmarshal(buf.Bytes(), &result); jsonErr != nil {
		t.Fatalf("invalid JSON: %v", jsonErr)
	}
	if result["ok"] != false {
		t.Error("expected ok=false")
	}
}

func TestHandleRootError_UnknownError(t *testing.T) {
	var buf bytes.Buffer
	code := handleRootError(&buf, errors.New("something unexpected"))
	if code != output.ExitInternal {
		t.Errorf("expected exit code %d (internal), got %d", output.ExitInternal, code)
	}
	var result map[string]any
	if jsonErr := json.Unmarshal(buf.Bytes(), &result); jsonErr != nil {
		t.Fatalf("invalid JSON: %v", jsonErr)
	}
	errObj := result["error"].(map[string]any)
	if errObj["code"] != "internal" {
		t.Errorf("expected error code 'internal', got %v", errObj["code"])
	}
}

func TestHandleRootError_ValidationError(t *testing.T) {
	var buf bytes.Buffer
	err := output.ErrValidation("Missing required argument: query")
	code := handleRootError(&buf, err)
	if code != output.ExitValidation {
		t.Errorf("expected exit code %d (validation), got %d", output.ExitValidation, code)
	}
}
