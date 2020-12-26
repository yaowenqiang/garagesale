package web

import (
    "github.com/go-chi/chi"
    "log"
    "syscall"
    "os"
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
    shutdown chan os.Signal
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

func NewApp(shutdown chan os.Signal, logger *log.Logger, mw ...Middleware) *App {
	app := App{
		Log: logger,
		mux: chi.NewRouter(),
		mw:  mw,
    shutdown: shutdown,
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
        //generate a shutdown
        //ctx = context.WithValue(ctx, KeyValues, &v)

        ctx = context.WithValue(ctx, KeyValues, &v)

        if err := h(ctx, w, r); err != nil {
            a.Log.Printf("%s: Unhandled error %+v", v.TraceID, err)
            if IsShutdown(err) {
                a.SignalShutdown()
            }
        }

    }
    a.mux.MethodFunc(method, pattern, fn)
}

// ServeHTTP implements the http.Handler interface.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.och.ServeHTTP(w, r)
}


// SignalShutdown is used to gracefully shutdown the app when an integrity
// issue is identified.
func (a *App) SignalShutdown() {
	a.Log.Println("error returned from handler indicated integrity issue, shutting down service")
	a.shutdown <- syscall.SIGSTOP
}
