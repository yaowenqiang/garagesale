package main
import (
    "net/http"
    "fmt"
    "log"
    "time"
    "context"
    "math/rand"
    "os"
    "os/signal"
    "syscall"
    "github.com/yaowenqiang/garagesale/internal/platform/conf"
    "github.com/yaowenqiang/garagesale/cmd/sales-api/internal/handlers"
    "github.com/yaowenqiang/garagesale/internal/platform/database"
)


func main() {
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
                log.Fatalf("error: generating config usage: %s", err)
            }
            fmt.Println(usage)
            return
        }
        log.Fatalf("error: parsing config: %s", err)
    }

    out, err := conf.String(&cfg)
    if err != nil {
        log.Fatalf("error : generating config for output : %s", err)
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
        log.Fatalf("error: connecting to db: %s", err)
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
        log.Fatalf("error: Listening and Serving %s", err)
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
            log.Printf("main: could not stop server gracefully %v", err)
        }
    }
    /*
    //h := http.HandlerFunc(Echo)
    h := http.HandlerFunc(Echo)
    log.Println("Listten on localhost:8111")
    if err := http.ListenAndServe("localhost:8111", h); err != nil {
        log.Fatal(err)
    }
    */
}

// the echo method

func Echo(w http.ResponseWriter, r *http.Request) {
    id := rand.Intn(300)
    fmt.Println("starting ", id)
    time.Sleep(3*time.Second)
    fmt.Fprintln(w,"You asked %s ", r.Method, r.URL.Path )
    fmt.Println("ending ", id)
}
