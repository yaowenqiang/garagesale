package main
import (
    "net/url"
    "log"
    "time"
    "flag"
    "github.com/yaowenqiang/garagesale/internal/schema"
    "github.com/yaowenqiang/garagesale/internal/platform/database"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Product struct {
    ID string `db:"product_id" json: "id"`
    Name string `db:"name" json: "name"`
    Cost int`db:"cost" json: "cost"`
    Quantity int`db:"quantity" json: "quantity"`
    DateCreated time.Time `db:"date_created" json: "date_created"`
    DateUpdated time.Time `db:"date_updated" json: "date_updated"`
}

func main() {
    log.Printf("main: Started")
    defer log.Printf("main: Completed")

    db, err := database.Open()

    if err != nil {
        log.Fatalf("error: connecting to db: %s", err)
    }

    defer db.Close()

    flag.Parse()
    switch flag.Arg(0) {
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
