package handlers

import (
    "log"
    "net/http"
	"github.com/jmoiron/sqlx"
    "github.com/yaowenqiang/garagesale/internal/platform/web"
	"github.com/yaowenqiang/garagesale/internal/platform/auth"
    "github.com/yaowenqiang/garagesale/internal/mid"
)

//handle all api routes

func API(logger *log.Logger, db *sqlx.DB, authenticator *auth.Authenticator) *web.App {
    app := web.NewApp(logger,mid.Logger(logger), mid.Errors(logger), mid.Metrics())
    p := Product {
        Db: db,
        Log: logger,
    }

    c := Check {
        DB: db,
    }

    u := Users{
        DB: db,
        authenticator: authenticator,
    }

	app.Handle(http.MethodGet, "/v1/users/token", u.Token)

	app.Handle(http.MethodGet, "/v1/health", c.Health)

    app.Handle(http.MethodGet, "/v1/products", p.List)
    app.Handle(http.MethodGet, "/v1/products/{id}", p.Retrieve)
    app.Handle(http.MethodPost, "/v1/products", p.Create)
	app.Handle(http.MethodPut, "/v1/products/{id}", p.Update)
	app.Handle(http.MethodDelete, "/v1/products/{id}", p.Delete)

	app.Handle(http.MethodPost, "/v1/products/{id}/sales", p.AddSale)
	app.Handle(http.MethodGet, "/v1/products/{id}/sales", p.ListSales)


    return app
}
