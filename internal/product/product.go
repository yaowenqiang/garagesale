package product
import (
	"github.com/jmoiron/sqlx"
)

func List(db *sqlx.DB) ([]Product, error) {
    //list products

    const q = "SELECT name, cost, quantity, date_updated, date_created  FROM products";
    list := []Product{}
    if err := db.Select(&list, q); err != nil {
        return nil, err
    }
    return list, nil

}

// get s single product
func Retrieve(db *sqlx.DB, id string) (*Product, error) {
    var p Product

    const q = "SELECT name, cost, quantity, date_updated, date_created  FROM products  where product_id = $1";

    if err := db.Get(&p, q, id); err != nil {
        return nil, err
    }
    return &p, nil

}
