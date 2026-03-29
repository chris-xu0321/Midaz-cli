package output

// Exit codes for the seer-q CLI.
//
// Fine-grained error types (not_found, unauthorized, etc.) are communicated
// via the JSON error envelope's "code" field, not via exit codes.
const (
	ExitOK         = 0 // Success
	ExitInternal   = 1 // Unexpected CLI error
	ExitValidation = 2 // Missing required arg, unknown flag/command
	ExitConfig     = 3 // Config file malformed or required key missing
	ExitNetwork    = 4 // Can't reach API, timeout
	ExitAPI        = 5 // HTTP 4xx/5xx from API
)
