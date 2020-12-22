package web

import (
    "github.com/go-chi/chi"
    "log"
    "net/http"
)



type App struct {
    mux *chi.Mux
    Log *log.Logger
    mw []Middleware
}


type Handler func(http.ResponseWriter, *http.Request) error

//new app

func NewApp(logger *log.Logger, mw ...Middleware) *App {
    return &App{
        mux: chi.NewRouter(),
        Log: logger,
        mw: mw,
    }
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    a.mux.ServeHTTP(w,r)
}


func (a *App) Handle(method, pattern string, h Handler) {
    h = wrapMiddleware(a.mw, h)
    fn := func(w http.ResponseWriter, r *http.Request) {
        if err := h(w, r); err != nil {
            a.Log.Printf("ERROR: Unhandled error %v", err)
        }
    }
    a.mux.MethodFunc(method, pattern, fn)
}
