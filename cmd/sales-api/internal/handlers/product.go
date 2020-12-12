package handlers

import (
    "encoding/json"
    "log"
    "net/http"

	"github.com/jmoiron/sqlx"
    "github.com/yaowenqiang/garagesale/internal/product"
)

type Product struct {
    Db *sqlx.DB
}

func (p *Product) List(w http.ResponseWriter, r *http.Request) {
    list, err := product.List(p.Db);
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Println("error query! db")
    }

    data, err := json.Marshal(list)

    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Println("err marshaling  ", err)
        return
    } else {
        w.Header().Set("Content-Type","application/json; charset=utf8")
        w.WriteHeader(http.StatusOK)
    }

    if _, err := w.Write(data); err != nil {
        log.Println("err writing ", err)
    }

}
