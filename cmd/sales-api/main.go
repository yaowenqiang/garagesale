package main
import (
    "net/http"
    "fmt"
    "log"
    "time"
    "context"
    "os"
    "os/signal"
    "syscall"
    "github.com/pkg/errors"
    "github.com/yaowenqiang/garagesale/internal/platform/conf"
    "github.com/yaowenqiang/garagesale/cmd/sales-api/internal/handlers"
    "github.com/yaowenqiang/garagesale/internal/platform/database"
)

func main() {
    if err := run(); err != nil {
        log.Fatal(err)
    }
}

func run() error {
    log.Printf("main: Started")
    defer log.Printf("main: Completed")


    var cfg struct {
		Web struct {
			Address         string        `conf:"default:localhost:8000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
        DB struct {
            User string `conf:"default:postgres"`
            Password string `conf:"default:postgres,noprint"`
            Host string `conf:"default:localhost"`
            Name string `conf:"default:postgres"`
            DisableTLS bool `conf:"default:false"`
        }
    }


    //parse configuration

    if err := conf.Parse(os.Args[1:], "sales", &cfg); err != nil {
        if err == conf.ErrHelpWanted {
            usage, err := conf.Usage("SALES", &cfg)
            if err != nil {
                return errors.Wrap(err, "generating config usage")
            }
            fmt.Println(usage)
            return nil
        }
        return errors.Wrap(err, "parsing config")
    }

    out, err := conf.String(&cfg)
    if err != nil {
        return errors.Wrap(err, "enerating config for output")
    }
    log.Printf("main: Config: \n%v\n", out)

    db, err := database.Open(database.Config{
        Host: cfg.DB.Host,
        User: cfg.DB.User,
        Name: cfg.DB.Name,
        Password: cfg.DB.Password,
        DisableTLS: cfg.DB.DisableTLS,
    })

    if err != nil {
        return errors.Wrap(err, "connecting to db")
    }

    defer db.Close()

    ps := handlers.Product{Db: db}

    api := http.Server{
        Addr: cfg.Web.Address,
        Handler: http.HandlerFunc(ps.List),
        ReadTimeout: cfg.Web.ReadTimeout,
        WriteTimeout: cfg.Web.WriteTimeout,
    }
    serverErrors := make(chan error, 1)

    go func() {
        log.Printf("main: API listen on %s", api.Addr)
        serverErrors <- api.ListenAndServe()
    }()

    shutdown := make(chan os.Signal, 1)
    signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

    select {
    case err := <- serverErrors:
        return errors.Wrap(err, "Listening and Serving")
    case <-shutdown:
        log.Println("main: Start shutdown")

        ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
        defer cancel()

        err := api.Shutdown(ctx)
        if err != nil {
            log.Printf("main: Graceful shutdown did not complete in %v : %v", cfg.Web.ShutdownTimeout, err)
            err = api.Close()
        }

        if err != nil {
            return errors.Wrap(err, "graceful shutdown")
        }
    }
    return nil
}

