package handlers

import (
    "log"
    "time"
    "net/http"

	"github.com/jmoiron/sqlx"
	"github.com/go-chi/chi"
    "github.com/yaowenqiang/garagesale/internal/product"
    "github.com/yaowenqiang/garagesale/internal/platform/web"
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

    if err:= web.Respond(w, list, http.StatusOK); err != nil {
        p.Log.Println("error responding", err)
        return
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

    if err:= web.Respond(w, prod, http.StatusOK); err != nil {
        p.Log.Println("error responding", err)
        return
    }

}

func (p *Product) Create(w http.ResponseWriter, r *http.Request) {
    var np product.NewProduct
    if err := web.Decode(r, &np); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        p.Log.Println(err)
        return
    }

    prod, err := product.Create(p.Db, np, time.Now())
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        p.Log.Println("error creating product", err)
    }

    if err := web.Respond(w, prod, http.StatusCreated); err != nil {
        p.Log.Println("error responding", err)
        return
    }


}
