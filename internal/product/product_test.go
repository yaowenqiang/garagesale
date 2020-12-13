package product

import (
    "testing"
    "time"
    "github.com/yaowenqiang/garagesale/internal/product"
    "github.com/yaowenqiang/garagesale/internal/schema"
    "github.com/yaowenqiang/garagesale/internal/platform/database/databasetest"
    "github.com/google/go-cmp/cmp"
)


func TestProducts(t *testing.T) {
    db, cleanup := databasetest.Setup(t)
    defer cleanup()
    np := product.NewProduct{
        Name: "new Comic book",
        Cost: 10,
        Quantity: 100,
    }

    now := time.Date(2020, time.January, 1, 0,0,0,0,time.UTC)
    p0, err := product.Create(db,np , now)
    if err != nil {
        t.Fatalf("could not create product %v", err)
    }

    p1, err := product.Retrieve(db, p0.ID)

    if err != nil {
        t.Fatalf("could not retrive product %v", err)
    }

    if diff := cmp.Diff(p0, p1); diff != "" {
        t.Fatalf("saved product did not match created see diff:\n%s", diff)
    }


}

func TestProductList(t *testing.T) {
    db, cleanup := databasetest.Setup(t)
    defer cleanup()

    if err := schema.Seed(db); err != nil {
        t.Fatal(err)
    }

    ps, err := product.List(db)

    if err != nil {
        t.Fatalf("Listing products: %s", err)
    }

    if exp, got := 2, len(ps); exp != got {
        t.Fatalf("expected product list size: %v got %v",exp, got )
    }

}
