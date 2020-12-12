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

func (p *Product) List(w http.ResponseWriter, r *http.Request) error {
    p.Log.Println("SALES")
    list, err := product.List(p.Db)
    if err != nil {
        return err
    }

    return  web.Respond(w, list, http.StatusOK)

}

//Retrieve single product
func (p *Product) Retrieve(w http.ResponseWriter, r *http.Request) error {
    p.Log.Println("SALES")
    id := chi.URLParam(r, "id")
    prod, err := product.Retrieve(p.Db, id)
    if err != nil {
        return err
    }

    return  web.Respond(w, prod, http.StatusOK)

}

func (p *Product) Create(w http.ResponseWriter, r *http.Request)  error {
    var np product.NewProduct
    if err := web.Decode(r, &np); err != nil {
        return err
    }

    prod, err := product.Create(p.Db, np, time.Now())
    if err != nil {
        return err
    }

    return  web.Respond(w, prod, http.StatusCreated)
}
