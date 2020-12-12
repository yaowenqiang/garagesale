package main
import (
    "net/url"
    "log"
    "os"
    "fmt"
    "github.com/yaowenqiang/garagesale/internal/schema"
    "github.com/yaowenqiang/garagesale/internal/platform/database"
    "github.com/yaowenqiang/garagesale/internal/platform/conf"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)


func main() {

	var cfg struct {
		DB struct {
			User       string `conf:"default:postgres"`
			Password   string `conf:"default:postgres,noprint"`
			Host       string `conf:"default:localhost"`
			Name       string `conf:"default:postgres"`
			DisableTLS bool   `conf:"default:false"`
		}
        Args conf.Args
	}
	if err := conf.Parse(os.Args[1:], "SALES", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("SALES", &cfg)
			if err != nil {
				log.Fatalf("error : generating config usage : %v", err)
			}
			fmt.Println(usage)
			return
		}
		log.Fatalf("error: parsing config: %s", err)
	}

    log.Printf("main: Started")
    defer log.Printf("main: Completed")

    db, err := database.Open(database.Config{
        User: cfg.DB.User,
        Password: cfg.DB.Password,
        Host: cfg.DB.Host,
        Name: cfg.DB.Name,
        DisableTLS: cfg.DB.DisableTLS,
    })

    if err != nil {
        log.Fatalf("error: connecting to db: %s", err)
    }

    defer db.Close()

    switch cfg.Args.Num(0) {
    case "migrate":
        if err := schema.Migrate(db); err != nil {
            log.Fatal("applying migrations ", err)
        }
        log.Println("Migration complete")
        return
    case "seed":
        if err := schema.Seed(db); err != nil {
            log.Fatal("applying seed data ", err)
        }
        log.Println("Seeding complete")
        return
    }

}


func openDB() (*sqlx.DB, error) {
    q := url.Values{}
    q.Set("sslmode", "disable")
    q.Set("timezone", "utc")

    u := url.URL{
        Scheme: "postgres",
        User: url.UserPassword("postgres", "postgres"),
        Host: "localhost",
        Path: "postgres",
        RawQuery: q.Encode(),
    }

    return sqlx.Open("postgres", u.String())
}
