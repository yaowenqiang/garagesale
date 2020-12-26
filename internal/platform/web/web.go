package web

import (
    "github.com/go-chi/chi"
    "log"
    "context"
    "time"
    "net/http"
	"go.opencensus.io/trace"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"
)



type App struct {
    mux *chi.Mux
    Log *log.Logger
    mw []Middleware
	och *ochttp.Handler
}


type Handler func(context.Context, http.ResponseWriter, *http.Request) error


type ctxKey int

const KeyValues ctxKey = 1

type Values struct {
    TraceID string
    StatusCode int
    Start time.Time
}

//new app

func NewApp(logger *log.Logger, mw ...Middleware) *App {
	app := App{
		Log: logger,
		mux: chi.NewRouter(),
		mw:  mw,
	}

	// Create an OpenCensus HTTP Handler which wraps the router. This will start
	// the initial span and annotate it with information about the request/response.
	//
	// This is configured to use the W3C TraceContext standard to set the remote
	// parent if an client request includes the appropriate headers.
	// https://w3c.github.io/trace-context/
	app.och = &ochttp.Handler{
		Handler:     app.mux,
		Propagation: &tracecontext.HTTPFormat{},
	}

	return &app
}



func (a *App) Handle(method, pattern string, h Handler, mw ...Middleware) {
	// First wrap handler specific middleware around this handler.
	h = wrapMiddleware(mw, h)

	// Add the application's general middleware to the handler chain.
	h = wrapMiddleware(a.mw, h)

    fn := func(w http.ResponseWriter, r *http.Request) {

        ctx , span := trace.StartSpan(r.Context(), "internal platform web")
        defer span.End()

        v := Values {
            TraceID: span.SpanContext().TraceID.String(),
            Start: time.Now(),
        }
        ctx = context.WithValue(ctx, KeyValues, &v)

        if err := h(ctx, w, r); err != nil {
            a.Log.Printf("ERROR: Unhandled error %v", err)
        }

    }
    a.mux.MethodFunc(method, pattern, fn)
}

// ServeHTTP implements the http.Handler interface.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.och.ServeHTTP(w, r)
}
