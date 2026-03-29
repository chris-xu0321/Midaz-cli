// Package output defines response envelopes, structured errors, exit codes,
// and writers for the seer-q CLI.
//
// Success envelopes go to stdout:
//
//	{ "ok": true, "data": ..., "meta": { ... } }
//
// Error envelopes go to stderr:
//
//	{ "ok": false, "error": { "code": "...", "message": "...", "hint": "..." } }
package output

// Envelope is the standard success response wrapper.
// Meta uses map[string]any (not a struct) to preserve zero-value ints
// like "contradicting_count": 0 which omitempty would drop.
type Envelope struct {
	OK   bool           `json:"ok"`
	Data interface{}    `json:"data"`
	Meta map[string]any `json:"meta"`
}

// ErrorEnvelope is the standard error response wrapper.
type ErrorEnvelope struct {
	OK    bool       `json:"ok"`
	Error *ErrDetail `json:"error"`
}

// ErrDetail describes a structured error.
type ErrDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Hint    string `json:"hint,omitempty"`
}
