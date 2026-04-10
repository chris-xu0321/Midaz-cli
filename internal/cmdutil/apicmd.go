package cmdutil

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/SparkssL/Midaz-cli/internal/client"
	"github.com/SparkssL/Midaz-cli/internal/output"
)

// NormalizeFn transforms raw API response bytes into envelope data and meta.
type NormalizeFn func(body []byte) (data interface{}, meta map[string]any, err error)

// APISpec describes an API command's HTTP call and response normalization.
type APISpec struct {
	Path      string
	Params    url.Values
	Normalize NormalizeFn
}

// RunAPICommand executes an API call and writes the result.
func RunAPICommand(f *Factory, opts *RunOpts, spec *APISpec) error {
	c, err := f.Client()
	if err != nil {
		return err
	}

	resp, err := c.Get(opts.Ctx, spec.Path, spec.Params)
	if err != nil {
		return err
	}

	if opts.Raw {
		return output.WriteRaw(opts.Out, resp.Body, opts.Format)
	}

	data, meta, err := spec.Normalize(resp.Body)
	if err != nil {
		return output.Errorf(output.ExitInternal, "internal", "failed to parse response: %s", err)
	}

	return output.WriteSuccess(opts.Out, data, meta, opts.Format)
}

// MutatingAPISpec describes a write (POST/PUT/DELETE) API call.
type MutatingAPISpec struct {
	Method    string // "POST", "PUT", "DELETE"
	Path      string
	Body      []byte      // JSON body (nil for DELETE)
	Normalize NormalizeFn // optional — if nil, uses NormalizePassthrough
}

// RunMutatingAPICommand executes a write API call and writes the result.
func RunMutatingAPICommand(f *Factory, opts *RunOpts, spec *MutatingAPISpec) error {
	c, err := f.Client()
	if err != nil {
		return err
	}

	var resp *client.Response
	switch spec.Method {
	case "POST":
		resp, err = c.Post(opts.Ctx, spec.Path, spec.Body)
	case "PUT":
		resp, err = c.Put(opts.Ctx, spec.Path, spec.Body)
	case "PATCH":
		resp, err = c.Patch(opts.Ctx, spec.Path, spec.Body)
	case "DELETE":
		resp, err = c.Delete(opts.Ctx, spec.Path)
	default:
		return output.Errorf(output.ExitInternal, "internal", "unknown method: %s", spec.Method)
	}
	if err != nil {
		return err
	}

	if opts.Raw {
		return output.WriteRaw(opts.Out, resp.Body, opts.Format)
	}

	normalize := spec.Normalize
	if normalize == nil {
		normalize = NormalizePassthrough
	}

	data, meta, err := normalize(resp.Body)
	if err != nil {
		return output.Errorf(output.ExitInternal, "internal", "failed to parse response: %s", err)
	}

	return output.WriteSuccess(opts.Out, data, meta, opts.Format)
}

// --- Shared normalizers ---

// NormalizeBareArray parses a JSON array and returns it with a count meta.
func NormalizeBareArray(body []byte) (interface{}, map[string]any, error) {
	var arr []json.RawMessage
	if err := json.Unmarshal(body, &arr); err != nil {
		return nil, nil, fmt.Errorf("expected JSON array: %w", err)
	}
	// Re-unmarshal to get proper interface{} slice for marshaling
	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, nil, err
	}
	return data, map[string]any{"count": len(arr)}, nil
}

// NormalizePassthrough returns the parsed JSON as-is with empty meta.
func NormalizePassthrough(body []byte) (interface{}, map[string]any, error) {
	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, nil, err
	}
	return data, map[string]any{}, nil
}

// --- Map helpers for custom normalizers ---

// ParseMap parses JSON into a map of raw messages, preserving all fields.
func ParseMap(body []byte) (map[string]json.RawMessage, error) {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, fmt.Errorf("expected JSON object: %w", err)
	}
	return m, nil
}

// RebuildMap converts a map[string]json.RawMessage back to a marshallable map.
func RebuildMap(m map[string]json.RawMessage) (interface{}, error) {
	result := make(map[string]interface{}, len(m))
	for k, v := range m {
		var val interface{}
		if err := json.Unmarshal(v, &val); err != nil {
			return nil, fmt.Errorf("failed to unmarshal key %q: %w", k, err)
		}
		result[k] = val
	}
	return result, nil
}

// CountArray counts elements in a JSON array. Returns 0 if not an array.
func CountArray(raw json.RawMessage) int {
	var arr []json.RawMessage
	if json.Unmarshal(raw, &arr) == nil {
		return len(arr)
	}
	return 0
}

// UnmarshalString extracts a Go string from a JSON string value.
func UnmarshalString(raw json.RawMessage) string {
	var s string
	json.Unmarshal(raw, &s)
	return s
}

// UnmarshalInt extracts a Go int from a JSON number value.
func UnmarshalInt(raw json.RawMessage) int {
	var n int
	json.Unmarshal(raw, &n)
	return n
}

// ExtractViewURL is a common helper that extracts view_url from a map into meta,
// deletes it from the map, and returns the url string.
func ExtractViewURL(m map[string]json.RawMessage) string {
	raw, ok := m["view_url"]
	if !ok {
		return ""
	}
	delete(m, "view_url")
	return UnmarshalString(raw)
}
