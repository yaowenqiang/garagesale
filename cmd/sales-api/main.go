package main
import (
    "net/http"
    _ "net/http/pprof" // register /debug/pprof handlers
    _ "expvar" //register  /debug/varshandler
	"crypto/rsa"
	"io/ioutil"
    "fmt"
    "log"
    "syscall"
    "time"
    "context"
    "os"
    "os/signal"
    "github.com/pkg/errors"
	"contrib.go.opencensus.io/exporter/zipkin"
    "github.com/yaowenqiang/garagesale/internal/platform/conf"
    "github.com/yaowenqiang/garagesale/cmd/sales-api/internal/handlers"
    "github.com/yaowenqiang/garagesale/internal/platform/database"
	"github.com/yaowenqiang/garagesale/internal/platform/auth"
	jwt "github.com/dgrijalva/jwt-go"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
	"go.opencensus.io/trace"
)

func main() {
    if err := run(); err != nil {
        log.Fatal(err)
    }
}

func run() error {

    log := log.New(os.Stdout, "SALES : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

    log.Printf("main: Started")
    defer log.Printf("main: Completed")


    var cfg struct {
		Web struct {
			Address         string        `conf:"default:localhost:8000"`
			Debug           string        `conf:"default:localhost:6060"`
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
		Auth struct {
			KeyID          string `conf:"default:1"`
			PrivateKeyFile string `conf:"default:private.pem"`
			Algorithm      string `conf:"default:RS256"`
		}
		Trace struct {
			URL         string  `conf:"default:http://localhost:9411/api/v2/spans"`
			Service     string  `conf:"default:sales-api"`
			Probability float64 `conf:"default:1"`
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

	// =========================================================================
	// Initialize authentication support

	authenticator, err := createAuth(
		cfg.Auth.PrivateKeyFile,
		cfg.Auth.KeyID,
		cfg.Auth.Algorithm,
	)
	if err != nil {
		return errors.Wrap(err, "constructing authenticator")
	}

	// =========================================================================

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

	// =========================================================================
	// Start Tracing Support

	closer, err := registerTracer(
		cfg.Trace.Service,
		cfg.Web.Address,
		cfg.Trace.URL,
		cfg.Trace.Probability,
	)
	if err != nil {
		return err
	}
	defer closer()

    //debug service
    go func() {
        log.Printf("main: Debug service listening on %s", cfg.Web.Debug)
        err := http.ListenAndServe(cfg.Web.Debug, http.DefaultServeMux)
        if err != nil {
            log.Printf("main: Debug Service ended %s", err)
        }
    }()


    shutdown := make(chan os.Signal, 1)
    signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

    api := http.Server{
        Addr: cfg.Web.Address,
        Handler: handlers.API(shutdown, log,db, authenticator),
        ReadTimeout: cfg.Web.ReadTimeout,
        WriteTimeout: cfg.Web.WriteTimeout,
    }
    serverErrors := make(chan error, 1)

    go func() {
        log.Printf("main: API listen on %s", api.Addr)
        serverErrors <- api.ListenAndServe()
    }()


    select {
    case err := <- serverErrors:
        return errors.Wrap(err, "Listening and Serving")
    case sig := <-shutdown:
        log.Println("main: Start shutdown", sig)

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

        if sig == syscall.SIGSTOP {
            return errors.New("integrity error detected, asking for self shutdown")
        }
    }
    return nil
}

func createAuth(privateKeyFile, keyID, algorithm string) (*auth.Authenticator, error) {

	keyContents, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "reading auth private key")
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyContents)
	if err != nil {
		return nil, errors.Wrap(err, "parsing auth private key")
	}

	public := auth.NewSimpleKeyLookupFunc(keyID, key.Public().(*rsa.PublicKey))

	return auth.NewAuthenticator(key, keyID, algorithm, public)
}

func registerTracer(service, httpAddr, traceURL string, probability float64) (func() error, error) {
	localEndpoint, err := openzipkin.NewEndpoint(service, httpAddr)
	if err != nil {
		return nil, errors.Wrap(err, "creating the local zipkinEndpoint")
	}
	reporter := zipkinHTTP.NewReporter(traceURL)

	trace.RegisterExporter(zipkin.NewExporter(reporter, localEndpoint))
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.ProbabilitySampler(probability),
	})

	return reporter.Close, nil
}
