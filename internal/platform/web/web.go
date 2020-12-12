package web

import (
    "github.com/go-chi/chi"
    "log"
    "net/http"
)



type App struct {
    mux *chi.Mux
    Log *log.Logger
}


//new app

func NewApp(logger *log.Logger) *App {
    return &App{
        mux: chi.NewRouter(),
        Log: logger,
    }
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    a.mux.ServeHTTP(w,r)
}


func (a *App) Handle(method, pattern string, fn http.HandlerFunc) {
    a.mux.MethodFunc(method, pattern, fn)
}
