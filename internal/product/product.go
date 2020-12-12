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
