package product
import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/google/uuid"
	"github.com/pkg/errors"
    "time"
    "context"
)

var (
    ErrNotFound  = errors.New("product not found")
    ErrInvalidID = errors.New("id provided was not a valid UUID")
)


func List(ctx context.Context, db *sqlx.DB) ([]Product, error) {
    //list products

    const q = "SELECT product_id, name, cost, quantity, date_updated, date_created  FROM products";
    list := []Product{}
    if err := db.SelectContext(ctx, &list, q); err != nil {
        return nil, err
    }
    return list, nil

}

// get s single product
func Retrieve(ctx context.Context, db *sqlx.DB, id string) (*Product, error) {
    var p Product

    if _, err := uuid.Parse(id); err != nil {
        return nil, ErrInvalidID
    }

    const q = "SELECT product_id, name, cost, quantity, date_updated, date_created  FROM products  where product_id = $1";

    if err := db.GetContext(ctx, &p, q, id); err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrNotFound
        }
        return nil, err
    }
    return &p, nil

}

func Create(ctx context.Context, db *sqlx.DB, np NewProduct, now time.Time) (*Product, error) {
    p := Product {
        ID:  uuid.New().String(),
        Name: np.Name,
        Cost: np.Cost,
        Quantity: np.Quantity,
        DateCreated: now.UTC(),
        DateUpdated: now.UTC(),
    }


    const q = "INSERT into products (product_id, name, cost, quantity, date_created, date_updated) VALUES ($1,$2,$3,$4,$5,$6) ";

    if _, err := db.ExecContext(ctx, q, p.ID, p.Name, p.Cost, p.Quantity, p.DateCreated, p.DateUpdated); err != nil {
        return nil, errors.Wrapf(err, "insert product %v", np)
    }

    return &p, nil
}
