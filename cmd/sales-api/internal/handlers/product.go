package handlers

import (
    "log"
    "time"
    "context"
    "net/http"

	"github.com/jmoiron/sqlx"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
    "github.com/yaowenqiang/garagesale/internal/product"
    "github.com/yaowenqiang/garagesale/internal/platform/web"
	"github.com/yaowenqiang/garagesale/internal/platform/auth"
	"go.opencensus.io/trace"
)

type Product struct {
    Db *sqlx.DB
    Log *log.Logger
}

func (p *Product) List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

    ctx , span := trace.StartSpan(ctx, "handles ProductList")
    defer span.End()
    p.Log.Println("SALES")
    list, err := product.List(ctx, p.Db)
    if err != nil {
        return err
    }

    return  web.Respond(ctx, w, list, http.StatusOK)

}

//Retrieve single product
func (p *Product) Retrieve(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
    p.Log.Println("SALES")
    id := chi.URLParam(r, "id")
    prod, err := product.Retrieve(ctx,p.Db, id)
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

    return  web.Respond(ctx, w, prod, http.StatusOK)

}

func (p *Product) Create(ctx context.Context, w http.ResponseWriter, r *http.Request)  error {
    claims, ok := ctx.Value(auth.Key).(auth.Claims)
    if !ok {
        return errors.New("auth claims not in context")
    }
    var np product.NewProduct
    if err := web.Decode(r, &np); err != nil {
        return err
    }

    prod, err := product.Create(ctx, p.Db, claims, np, time.Now())
    if err != nil {
        return err
    }

    return  web.Respond(ctx, w, prod, http.StatusCreated)
}

// AddSale creates a new Sale for a particular product. It looks for a JSON
// object in the request body. The full model is returned to the caller.
func (p *Product) AddSale(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var ns product.NewSale
	if err := web.Decode(r, &ns); err != nil {
		return errors.Wrap(err, "decoding new sale")
	}

	productID := chi.URLParam(r, "id")

	sale, err := product.AddSale(ctx, p.Db, ns, productID, time.Now())
	if err != nil {
		return errors.Wrap(err, "adding new sale")
	}

	return web.Respond(ctx, w, sale, http.StatusCreated)
}

// ListSales gets all sales for a particular product.
func (p *Product) ListSales(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	list, err := product.ListSales(ctx, p.Db, id)
	if err != nil {
		return errors.Wrap(err, "getting sales list")
	}

	return web.Respond(ctx, w, list, http.StatusOK)
}

// Update decodes the body of a request to update an existing product. The ID
// of the product is part of the request URL.
func (p *Product) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

    claims, ok := ctx.Value(auth.Key).(auth.Claims)
    if !ok {
        return errors.New("auth claims not in context")
    }
	id := chi.URLParam(r, "id")

	var update product.UpdateProduct
	if err := web.Decode(r, &update); err != nil {
		return errors.Wrap(err, "decoding product update")
	}

	if err := product.Update(ctx, p.Db, claims, id, update, time.Now()); err != nil {
		switch err {
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case product.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "updating product %q", id)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (p *Product)Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
    if err := product.Delete(ctx, p.Db, id); err != nil {
        switch err {
        case product.ErrInvalidID:
            return web.NewRequestError(err, http.StatusBadRequest)
        default:
            return errors.Wrapf(err, "deleting product %s", id)
        }
    }
    return web.Respond(ctx, w, nil, http.StatusNoContent)
}
