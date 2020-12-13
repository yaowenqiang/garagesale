package tests
import (
    "encoding/json"
    _ "fmt"
    "log"
    "net/http"
    "net/http/httptest"
    "os"
    "strings"
    "testing"

    "github.com/yaowenqiang/garagesale/cmd/sales-api/internal/handlers"
    "github.com/yaowenqiang/garagesale/internal/platform/database/databasetest"
    "github.com/yaowenqiang/garagesale/internal/schema"
    "github.com/google/go-cmp/cmp"
)


func TestProducts(t *testing.T) {
    db, teardown := databasetest.Setup(t)
    defer teardown()

    if err := schema.Seed(db); err != nil {
        t.Fatal(err)
    }

    log := log.New(os.Stderr, "TEST: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

    tests := ProductTests{app: handlers.API(log,db)}
    t.Run("List", tests.List)
    t.Run("ProductCRUD", tests.ProductCRUD)
}

type ProductTests struct {
    app  http.Handler
}

func (p *ProductTests) List(t *testing.T) {
    req := httptest.NewRequest("GET", "/v1/products", nil)
    resp := httptest.NewRecorder()

    p.app.ServeHTTP(resp, req)

    if resp.Code != http.StatusOK {
        t.Fatalf("getting: expected status code %v, got %v", http.StatusOK, resp.Code)
    }

    var list []map[string]interface{}

    if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
        t.Fatalf("decoding: %s", err)
    }

	want := []map[string]interface{}{
		{
			"id":           "a2b0639f-2cc6-44b8-b97b-15d69dbb511e",
			"name":         "Comic Books",
			"cost":         float64(50),
			"quantity":     float64(42),
			"date_created": "2019-01-01T00:00:01.000001Z",
			"date_updated": "2019-01-01T00:00:01.000001Z",
		},
		{
			"id":           "72f8b983-3eb4-48db-9ed0-e45cc6bd716b",
			"name":         "McDonalds Toys",
			"cost":         float64(75),
			"quantity":     float64(120),
			"date_created": "2019-01-01T00:00:02.000001Z",
			"date_updated": "2019-01-01T00:00:02.000001Z",
		},
	}

    if diff := cmp.Diff(want, list); diff != "" {
        t.Fatalf("Response did not match expected, Diff\n%s", diff)
    }

}


func (p *ProductTests)ProductCRUD(t *testing.T) {
    {
        var created map[string]interface{}
		body := strings.NewReader(`{"name":"product0","cost":55,"quantity":6}`)
        req := httptest.NewRequest("POST", "/v1/products", body)

        req.Header.Set("Content-Type","application/json")

        resp := httptest.NewRecorder()

        p.app.ServeHTTP(resp, req)

        if http.StatusCreated != resp.Code {
            t.Fatalf("posting: expected status code %v, got %v", http.StatusCreated, resp.Code)
        }

        if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
            t.Fatalf("decoding %s", err)
        }

        if created["id"] == "" || created["id"] == nil {
            t.Fatal("expected non-empty product id")
        }

        if created["date_created"] == "" || created["date_created"] == nil {
            t.Fatal("expected non-empty product date_created")
        }

        if created["date_updated"] == "" || created["date_updated"] == nil {
            t.Fatal("expected non-empty product date_updated")
        }

		want := map[string]interface{}{
			"id":           created["id"],
			"date_created": created["date_created"],
			"date_updated": created["date_updated"],
			"name":         "product0",
			"cost":         float64(55),
			"quantity":     float64(6),
		}

        if diff := cmp.Diff(want, created); diff != "" {
            t.Fatalf("Response did not match expected. Diff:\n%s", diff)
        }
    }
}
