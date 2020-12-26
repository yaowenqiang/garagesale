package web

import (
    "net/http"
    "context"
    "encoding/json"
    "github.com/pkg/errors"
)


func Respond(ctx context.Context, w http.ResponseWriter, val interface{}, statusCode int) error {

    v, ok := ctx.Value(KeyValues).(*Values)

    if !ok {
        return errors.New("web values missing from context")
    }

    v.StatusCode = statusCode

    if statusCode == http.StatusNoContent {
        w.WriteHeader(statusCode)
        return nil
    }
    data, err := json.Marshal(val)

    if err != nil {
        return errors.Wrap(err, "Mashaling value to json")
    }
    w.Header().Set("Content-Type","application/json; charset=utf8")
    w.WriteHeader(statusCode)

    if _, err := w.Write(data); err != nil {
        return errors.Wrap(err, "Writing to client")
    }

    return nil
}

// RespondError sends an error reponse back to the client.
func RespondError(ctx context.Context, w http.ResponseWriter, err error) error {

	// If the error was of the type *Error, the handler has
	// a specific status code and error to return.
	if webErr, ok := errors.Cause(err).(*Error); ok {
		er := ErrorResponse{
			Error:  webErr.Err.Error(),
			Fields: webErr.Fields,
		}
		if err := Respond(ctx, w, er, webErr.Status); err != nil {
			return err
		}
		return nil
	}

	// If not, the handler sent any arbitrary error value so use 500.
	er := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}
	if err := Respond(ctx, w, er, http.StatusInternalServerError); err != nil {
		return err
	}
	return nil
}
