package web

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

