package web

import (
    "net/http"
    "encoding/json"
    "github.com/pkg/errors"
)


func Respond(w http.ResponseWriter, val interface{}, statusCode int) error {

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

func RespondError(w http.ResponseWriter, err error) error {
    if webErr, ok := errors.Cause(err).(*Error); ok {
        er := ErrorResponse{
            Error: webErr.Err.Error(),
        }

        if err :=  Respond(w, er, webErr.Status); err != nil {
            return err
        }
        return nil
    }

    er := ErrorResponse{
        Error: http.StatusText(http.StatusInternalServerError),
    }

    if err :=  Respond(w, er, http.StatusInternalServerError) ; err != nil {
        return err
    }

    return nil

}
