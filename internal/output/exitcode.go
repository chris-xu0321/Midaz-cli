package output

// Exit codes for the midaz CLI.
//
// Fine-grained error types (not_found, unauthorized, etc.) are communicated
// via the JSON error envelope's "code" field, not via exit codes.
const (
	ExitOK           = 0 // Success
	ExitInternal     = 1 // Unexpected CLI error
	ExitValidation   = 2 // Missing required arg, unknown flag/command
	ExitConfig       = 3 // Config file malformed or required key missing
	ExitNetwork      = 4 // Can't reach API, timeout
	ExitAPI          = 5 // HTTP 4xx/5xx from API
	ExitAuth         = 6 // Not authenticated (401) or credentials invalid
	ExitSubscription = 7 // Subscription required (402) or workspace paused
)
