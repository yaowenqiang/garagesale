package web

type ErrorResponse struct {
    Error string `json:"error"`
}

type Error struct {
    Err error
    Status int
}

func NewRequestError(err error, status int) error {
    return &Error{err, status}
}

func (e *Error) Error() string{
    return e.Err.Error()
}
