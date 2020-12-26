package web

import (
	"github.com/pkg/errors"
)

// ErrorResponse is the form used for API responses from failures in the API.
type ErrorResponse struct {
	Error  string       `json:"error"`
	Fields []FieldError `json:"fields,omitempty"`
}

type Error struct {
    Err error
    Status int
	Fields []FieldError
}

func NewRequestError(err error, status int) error {
    return &Error{err, status, nil}
}

func (e *Error) Error() string{
    return e.Err.Error()
}
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// shutdown is a type used to help with the graceful termination of the service.
type shutdown struct {
	Message string
}

// Error is the implementation of the error interface.
func (s *shutdown) Error() string {
	return s.Message
}

// NewShutdownError returns an error that causes the framework to signal
// a graceful shutdown.
func NewShutdownError(message string) error {
	return &shutdown{message}
}

// IsShutdown checks to see if the shutdown error is contained
// in the specified error value.
func IsShutdown(err error) bool {
	if _, ok := errors.Cause(err).(*shutdown); ok {
		return true
	}
	return false
}
