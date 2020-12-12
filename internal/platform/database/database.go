package database
import (
    "net/url"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // register ahe postgres database sql driver
)

//db connection config
type Config struct {
    Host string
    User string
    Name string
    Password string
    DisableTLS bool
}

func Open(cfg Config) (*sqlx.DB, error) {
    q := url.Values{}
    q.Set("sslmode", "require")
    if cfg.DisableTLS {
    q.Set("sslmode", "disable")
    }
    q.Set("timezone", "utc")

    u := url.URL{
        Scheme: "postgres",
        User: url.UserPassword(cfg.User, cfg.Password),
        Host: cfg.Host,
        Path: cfg.Name,
        RawQuery: q.Encode(),
    }

    return sqlx.Open("postgres", u.String())
}
