package web

import (
    "net/http"
    "encoding/json"
    "github.com/pkg/errors"
)


func Decode(r *http.Request, val interface{}) error {

    if err := json.NewDecoder(r.Body).Decode(&val); err != nil {
        return errors.Wrap(err, "decoding request body")
    }
    return nil

}
