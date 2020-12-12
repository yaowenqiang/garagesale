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
