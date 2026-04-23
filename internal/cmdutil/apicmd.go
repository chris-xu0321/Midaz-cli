package cmdutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/SparkssL/Midaz-cli/internal/client"
	"github.com/SparkssL/Midaz-cli/internal/output"
)

// NormalizeFn transforms raw API response bytes into envelope data and meta.
type NormalizeFn func(body []byte) (data interface{}, meta map[string]any, err error)

// APISpec describes an API command's HTTP call and response normalization.
//
// Method defaults to GET when empty. For POST/PATCH, set Body to a value that
// serializes to JSON (nil = empty body). Params is used for GET query strings
// only (ignored for other methods).
type APISpec struct {
	Method    string
	Path      string
	Params    url.Values
	Body      any
	Normalize NormalizeFn
}

// RunAPICommand executes an API call and writes the result.
// If Normalize is nil, falls back to NormalizePassthrough.
func RunAPICommand(f *Factory, opts *RunOpts, spec *APISpec) error {
	c, err := f.Client()
	if err != nil {
		return err
	}

	resp, err := callAPI(c, opts, spec)
	if err != nil {
		return err
	}

	if opts.Raw {
		return output.WriteRaw(opts.Out, resp.Body, opts.Format)
	}

	norm := spec.Normalize
	if norm == nil {
		norm = NormalizePassthrough
	}
	data, meta, err := norm(resp.Body)
	if err != nil {
		return output.Errorf(output.ExitInternal, "internal", "failed to parse response: %s", err)
	}

	return output.WriteSuccess(opts.Out, data, meta, opts.Format)
}

func callAPI(c *client.Client, opts *RunOpts, spec *APISpec) (*client.Response, error) {
	method := spec.Method
	if method == "" {
		method = http.MethodGet
	}
	switch method {
	case http.MethodGet:
		return c.Get(opts.Ctx, spec.Path, spec.Params)
	case http.MethodPost:
		return c.Post(opts.Ctx, spec.Path, spec.Body)
	case http.MethodPatch:
		return c.Patch(opts.Ctx, spec.Path, spec.Body)
	case http.MethodDelete:
		return c.Delete(opts.Ctx, spec.Path, spec.Body)
	default:
		return nil, output.Errorf(output.ExitInternal, "internal", "unsupported HTTP method: %s", method)
	}
}

// --- Shared normalizers ---

// NormalizeBareArray parses a JSON array and returns it with a count meta.
func NormalizeBareArray(body []byte) (interface{}, map[string]any, error) {
	var arr []json.RawMessage
	if err := json.Unmarshal(body, &arr); err != nil {
		return nil, nil, fmt.Errorf("expected JSON array: %w", err)
	}
	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, nil, err
	}
	return data, map[string]any{"count": len(arr)}, nil
}

// InjectItemViewURL returns a NormalizeFn for JSON-array responses that adds
// a `view_url` field to each object in the array, built from the item's `id`
// via pattern(id). Items that already have a `view_url` are left untouched,
// as are items without an `id` or when pattern returns "".
//
// Used to synthesize per-entity links on list endpoints (e.g. `midaz drivers`,
// `midaz theses`) where the upstream API only returns a bare array. Skills
// then have a URL to attach to every named thesis/driver in bulk replies.
func InjectItemViewURL(pattern func(id string) string) NormalizeFn {
	return func(body []byte) (interface{}, map[string]any, error) {
		var arr []map[string]json.RawMessage
		if err := json.Unmarshal(body, &arr); err != nil {
			return nil, nil, fmt.Errorf("expected JSON array of objects: %w", err)
		}
		for _, item := range arr {
			if _, has := item["view_url"]; has {
				continue
			}
			id := UnmarshalString(item["id"])
			if id == "" {
				continue
			}
			u := pattern(id)
			if u == "" {
				continue
			}
			raw, err := json.Marshal(u)
			if err != nil {
				continue
			}
			item["view_url"] = raw
		}
		data := make([]interface{}, len(arr))
		for i, item := range arr {
			rebuilt, err := RebuildMap(item)
			if err != nil {
				return nil, nil, err
			}
			data[i] = rebuilt
		}
		return data, map[string]any{"count": len(arr)}, nil
	}
}

// NormalizePassthrough returns the parsed JSON as-is with empty meta.
// Also accepts empty bodies (returns nil data).
func NormalizePassthrough(body []byte) (interface{}, map[string]any, error) {
	if len(body) == 0 {
		return nil, map[string]any{}, nil
	}
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
	_ = json.Unmarshal(raw, &s)
	return s
}

// UnmarshalInt extracts a Go int from a JSON number value.
func UnmarshalInt(raw json.RawMessage) int {
	var n int
	_ = json.Unmarshal(raw, &n)
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
