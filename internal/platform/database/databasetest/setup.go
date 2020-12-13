package databasetest
import (
    "testing"
    "time"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // register ahe postgres database sql driver
    "github.com/yaowenqiang/garagesale/internal/schema"
    "github.com/yaowenqiang/garagesale/internal/platform/database"
)


func Setup(t *testing.T) (*sqlx.DB, func()) {
    t.Helper()

    c := StartContainer(t)

    db, err := database.Open(database.Config{
        User: "postgres",
        Password: "postgres",
        Host: c.Host,
        Name: "postgres",
        DisableTLS: true,
    })

    if err != nil {
        t.Fatalf("opening database connection: %v", err)
    }
    t.Log("waiting for database to be ready")

    var pingError error

    maxAttempts := 20

    for attempts := 1; attempts <= maxAttempts; attempts++ {
        pingError := db.Ping()
        if pingError == nil {
            break
        }

        time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
    }

    if pingError != nil {
        StopContainer(t, c)
        t.Fatalf("waiting for database to be ready %v", pingError)
    }

    if err := schema.Migrate(db); err != nil {
        StopContainer(t, c)
        t.Fatalf("migrating %s", err)
    }

    teardown := func() {
        t.Helper()
        db.Close()
        StopContainer(t,c)
    }

    return db, teardown
}

