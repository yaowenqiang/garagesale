package handlers

import (
    "log"
    "time"
    "net/http"

	"github.com/jmoiron/sqlx"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
    "github.com/yaowenqiang/garagesale/internal/product"
    "github.com/yaowenqiang/garagesale/internal/platform/web"
)

type Product struct {
    Db *sqlx.DB
    Log *log.Logger
}

func (p *Product) List(w http.ResponseWriter, r *http.Request) error {
    p.Log.Println("SALES")
    list, err := product.List(r.Context(), p.Db)
    if err != nil {
        return err
    }

    return  web.Respond(w, list, http.StatusOK)

}

//Retrieve single product
func (p *Product) Retrieve(w http.ResponseWriter, r *http.Request) error {
    p.Log.Println("SALES")
    id := chi.URLParam(r, "id")
    prod, err := product.Retrieve(r.Context(),p.Db, id)
    if err != nil {
        switch err {
        case product.ErrNotFound:
            return web.NewRequestError(err, http.StatusNotFound)
        case product.ErrInvalidID:
            return web.NewRequestError(err, http.StatusBadRequest)
        default:
            return errors.Wrapf(err, "Looking for product %q", id)
        }

    }

    return  web.Respond(w, prod, http.StatusOK)

}

func (p *Product) Create(w http.ResponseWriter, r *http.Request)  error {
    var np product.NewProduct
    if err := web.Decode(r, &np); err != nil {
        return err
    }

    prod, err := product.Create(r.Context(), p.Db, np, time.Now())
    if err != nil {
        return err
    }

    return  web.Respond(w, prod, http.StatusCreated)
}
