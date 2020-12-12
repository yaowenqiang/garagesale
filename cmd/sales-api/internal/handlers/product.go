package handlers

import (
    "encoding/json"
    "log"
    "time"
    "net/http"

	"github.com/jmoiron/sqlx"
	"github.com/go-chi/chi"
    "github.com/yaowenqiang/garagesale/internal/product"
)

type Product struct {
    Db *sqlx.DB
    Log *log.Logger
}

func (p *Product) List(w http.ResponseWriter, r *http.Request) {
    p.Log.Println("SALES")
    list, err := product.List(p.Db)
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

//Retrieve single product
func (p *Product) Retrieve(w http.ResponseWriter, r *http.Request) {
    p.Log.Println("SALES")
    id := chi.URLParam(r, "id")
    prod, err := product.Retrieve(p.Db, id)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        p.Log.Println("error query! db")
    }

    data, err := json.Marshal(prod)

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

func (p *Product) Create(w http.ResponseWriter, r *http.Request) {
    var np product.NewProduct
    if err := json.NewDecoder(r.Body).Decode(&np); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        p.Log.Println(err)
        return
    }

    prod, err := product.Create(p.Db, np, time.Now())
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        p.Log.Println("error creating product", err)
    }

    data, err := json.Marshal(prod)

    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        p.Log.Println("err marshaling  ", err)
        return
    } else {
        w.Header().Set("Content-Type","application/json; charset=utf8")
        w.WriteHeader(http.StatusCreated)
    }

    if _, err := w.Write(data); err != nil {
        p.Log.Println("err writing ", err)
    }


}
