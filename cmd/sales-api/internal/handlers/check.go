package handlers

import (
    "net/http"

	"github.com/jmoiron/sqlx"
    "github.com/yaowenqiang/garagesale/internal/platform/web"
    "github.com/yaowenqiang/garagesale/internal/platform/database"
)
type Check struct {
    DB *sqlx.DB
}


func (c *Check) Health(w http.ResponseWriter, r *http.Request) error {
    var health struct {
        Status string `json.status`
    }

    if err := database.StatusCheck(r.Context(), c.DB); err != nil {
        health.Status = "db not ready"
        return web.Respond(r.Context(), w, health, http.StatusInternalServerError)
    }
    health.Status = "OK"
    return web.Respond(r.Context(), w, health, http.StatusOK)
}

