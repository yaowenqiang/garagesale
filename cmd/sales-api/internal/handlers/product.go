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
    Log *log.Logger
}

func (p *Product) List(w http.ResponseWriter, r *http.Request) {
    p.Log.Println("SALES")
    list, err := product.List(p.Db);
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        p.Log.Println("error query! db")
    }

    data, err := json.Marshal(list)

    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        p.Log.Println("err marshaling  ", err)
        return
    } else {
        w.Header().Set("Content-Type","application/json; charset=utf8")
        w.WriteHeader(http.StatusOK)
    }

    if _, err := w.Write(data); err != nil {
        p.Log.Println("err writing ", err)
    }

}
