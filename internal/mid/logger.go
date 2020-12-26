package mid

import (
	"log"
	"time"
	"context"
	"net/http"
    "github.com/pkg/errors"

	"github.com/yaowenqiang/garagesale/internal/platform/web"
)

// Logger will log a line for every request
func Logger(log *log.Logger) web.Middleware {

	// This is the actual middleware function to be executed.
	f := func(before web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
            v, ok := ctx.Value(web.KeyValues).(*web.Values)
            if !ok {
                return errors.New("web values missing from context")
            }


			// Run the handler chain and catch any propagated error.
            err := before(ctx, w, r)
            log.Printf(
                "%d %s %s (%v)",
                v.StatusCode,
                r.Method,
                r.URL.Path,
                time.Since(v.Start),
            )

                return err

		}

		return h
	}

	return f
}

