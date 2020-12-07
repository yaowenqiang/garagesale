package main
import (
    "net/http"
    "fmt"
    "log"
)


func main() {
    h := http.HandlerFunc(Echo)
    log.Println("Listten on localhost:8111")
    if err := http.ListenAndServe("localhost:8111", h); err != nil {
        log.Fatal(err)
    }
}

// the echo method

func Echo(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w,"You asked %s ", r.Method, r.URL.Path )
}

