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
)


func main() {
    log.Printf("main: Started")
    defer log.Printf("main: Completed")

    api := http.Server{
        Addr: "localhost:8111",
        Handler: http.HandlerFunc(Echo),
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

