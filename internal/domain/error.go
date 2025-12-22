package domain

import (
	"errors"
	"fmt"
)

// Error codes for mapping to HTTP status
const (
	EINVALID      = "invalid"       // 400
	EUNAUTHORIZED = "unauthorized"  // 401
	EFORBIDDEN    = "forbidden"     // 403
	ENOTFOUND     = "not_found"     // 404
	ECONFLICT     = "conflict"      // 409
	EINTERNAL     = "internal"      // 500
)

// Error represents a domain error with code, operation, and message.
type Error struct {
	Code    string
	Op      string
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Op, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Op, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Err
}

// Errorf creates a new domain error.
func Errorf(code, op, message string, args ...interface{}) *Error {
	return &Error{
		Code:    code,
		Op:      op,
		Message: fmt.Sprintf(message, args...),
	}
}

// WrapError wraps an existing error with domain context.
func WrapError(code, op, message string, err error) *Error {
	return &Error{
		Code:    code,
		Op:      op,
		Message: message,
		Err:     err,
	}
}

// ErrorCode returns the error code from a domain error.
func ErrorCode(err error) string {
	if err == nil {
		return ""
	}
	var e *Error
	if errors.As(err, &e) {
		return e.Code
	}
	return EINTERNAL
}

// ErrorMessage returns a user-friendly message from an error.
func ErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	var e *Error
	if errors.As(err, &e) {
		return e.Message
	}
	return "An unexpected error occurred"
}
