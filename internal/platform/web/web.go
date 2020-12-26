package web

import (
    "github.com/go-chi/chi"
    "log"
    "context"
    "time"
    "net/http"
)



type App struct {
    mux *chi.Mux
    Log *log.Logger
    mw []Middleware
}


type Handler func(http.ResponseWriter, *http.Request) error


type ctxKey int

const KeyValues ctxKey = 1

type Values struct {
    StatusCode int
    Start time.Time
}

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
        v := Values {
            Start: time.Now(),
        }
        ctx := r.Context()
        ctx = context.WithValue(ctx, KeyValues, &v)
        r = r.WithContext(ctx)

        if err := h(w, r); err != nil {
            a.Log.Printf("ERROR: Unhandled error %v", err)
        }
    }
    a.mux.MethodFunc(method, pattern, fn)
}
