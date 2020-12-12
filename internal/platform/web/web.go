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


type handler func(http.ResponseWriter, *http.Request) error

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


func (a *App) Handle(method, pattern string, h handler) {
    fn := func(w http.ResponseWriter, r *http.Request) {
        if err := h(w, r); err != nil {
            resp := ErrorResponse{
                Error: err.Error(),
            }

            if err:= Respond(w, resp, http.StatusInternalServerError); err != nil {
                a.Log.Println(err)
            }

        }
    }
    a.mux.MethodFunc(method, pattern, fn)
}
