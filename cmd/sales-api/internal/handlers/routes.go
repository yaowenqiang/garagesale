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
    app := web.NewApp(logger,mid.Logger(logger), mid.Errors(logger), mid.Metrics(), mid.Panics(logger))
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

	app.Handle(http.MethodGet, "/v1/health", c.Health, mid.Authenticate(authenticator))

    app.Handle(http.MethodGet, "/v1/products", p.List, mid.Authenticate(authenticator))

    app.Handle(http.MethodGet, "/v1/products/{id}", p.Retrieve, mid.Authenticate(authenticator))
    app.Handle(http.MethodPost, "/v1/products", p.Create, mid.Authenticate(authenticator))
	app.Handle(http.MethodPut, "/v1/products/{id}", p.Update, mid.Authenticate(authenticator))
	app.Handle(http.MethodDelete, "/v1/products/{id}", p.Delete, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))

	app.Handle(http.MethodPost, "/v1/products/{id}/sales", p.AddSale, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle(http.MethodGet, "/v1/products/{id}/sales", p.ListSales, mid.Authenticate(authenticator))


    return app
}
