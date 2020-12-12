package handlers

import (
    "log"
    "net/http"
	"github.com/jmoiron/sqlx"
    "github.com/yaowenqiang/garagesale/internal/platform/web"
)

//handle all api routes

func API(logger *log.Logger, db *sqlx.DB) *web.App {
    app := web.NewApp(logger)
    p := Product {
        Db: db,
        Log: logger,
    }

    app.Handle(http.MethodGet, "/v1/products", p.List)
    app.Handle(http.MethodGet, "/v1/products/{id}", p.Retrieve)
    return app
}
