package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// ExitError is a structured error that carries an exit code and optional detail.
// It is propagated up the call chain and handled by the root command to produce
// a JSON error envelope on stderr and the correct exit code.
type ExitError struct {
	Code   int
	Detail *ErrDetail
	Err    error
}

func (e *ExitError) Error() string {
	if e.Detail != nil {
		return e.Detail.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return fmt.Sprintf("exit %d", e.Code)
}

func (e *ExitError) Unwrap() error {
	return e.Err
}

// WriteErrorEnvelope writes a JSON error envelope to w.
func WriteErrorEnvelope(w io.Writer, err *ExitError) {
	if err.Detail == nil {
		return
	}
	env := ErrorEnvelope{
		OK:    false,
		Error: err.Detail,
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if encErr := enc.Encode(env); encErr != nil {
		return
	}
	buf.WriteTo(w)
}

// --- Convenience constructors ---

// Errorf creates an ExitError with the given exit code, error code, and formatted message.
func Errorf(exitCode int, errCode, format string, args ...any) *ExitError {
	var err error
	for _, arg := range args {
		if e, ok := arg.(error); ok {
			err = e
			break
		}
	}
	return &ExitError{
		Code:   exitCode,
		Detail: &ErrDetail{Code: errCode, Message: fmt.Sprintf(format, args...)},
		Err:    err,
	}
}

// ErrValidation creates a validation ExitError (exit 2).
func ErrValidation(format string, args ...any) *ExitError {
	return Errorf(ExitValidation, "validation", format, args...)
}

// ErrConfig creates a config ExitError (exit 3).
func ErrConfig(format string, args ...any) *ExitError {
	return Errorf(ExitConfig, "config", format, args...)
}

// ErrNetwork creates a network ExitError (exit 4).
func ErrNetwork(format string, args ...any) *ExitError {
	return Errorf(ExitNetwork, "network", format, args...)
}

// ErrAPI creates an API ExitError (exit 5) with the given error code.
func ErrAPI(errCode, format string, args ...any) *ExitError {
	return Errorf(ExitAPI, errCode, format, args...)
}

// ErrWithHint creates an ExitError with a hint string.
func ErrWithHint(exitCode int, errCode, msg, hint string) *ExitError {
	return &ExitError{
		Code:   exitCode,
		Detail: &ErrDetail{Code: errCode, Message: msg, Hint: hint},
	}
}
