package mid

import (
	"log"
	"time"
	"net/http"

	"github.com/yaowenqiang/garagesale/internal/platform/web"
)

// Logger will log a line for every request
func Logger(log *log.Logger) web.Middleware {

	// This is the actual middleware function to be executed.
	f := func(before web.Handler) web.Handler {

		h := func(w http.ResponseWriter, r *http.Request) error {
            now := time.Now()


			// Run the handler chain and catch any propagated error.
            err := before(w, r)
            log.Printf(
                "%s %s %v",
                r.Method,
                r.URL.Path,
                time.Since(now),
            )

                return err

		}

		return h
	}

	return f
}

