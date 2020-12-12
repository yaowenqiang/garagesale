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
    "encoding/json"
    "github.com/yaowenqiang/garagesale/internal/platform/database"
    "github.com/yaowenqiang/garagesale/cmd/sales-api/internal/handlers"
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

    ps := handlers.Product{Db: db}

    api := http.Server{
        Addr: "localhost:8111",
        Handler: http.HandlerFunc(ps.List),
        ReadTimeout: 5 * time.Second,
        WriteTimeout: 5 * time.Second,
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

        const timeout = 5 * time.Second
        ctx, cancel := context.WithTimeout(context.Background(), timeout)
        defer cancel()

        err := api.Shutdown(ctx)
        if err != nil {
            log.Printf("main: Graceful shutdown did not complete in %v : %v", timeout, err)
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


func ListProducts(w http.ResponseWriter, r *http.Request) {
    list := []Product{
        { Name: "comic books",Cost: 100, Quantity: 10,},
        { Name: "it books",Cost: 500, Quantity: 100},
    }


    data, err := json.Marshal(list)

    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.Println("err marshaling  ", err)
        return
    } else {
        w.Header().Set("Content-Type","application/json; charset=utf8")
        w.WriteHeader(http.StatusOK)
    }

    if _, err := w.Write(data); err != nil {
        log.Println("err writing ", err)
    }
}
