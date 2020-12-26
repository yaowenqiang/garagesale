package product
import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/google/uuid"
	"github.com/pkg/errors"
    "time"
    "context"
	"github.com/yaowenqiang/garagesale/internal/platform/auth"
)

var (
    ErrNotFound  = errors.New("product not found")
    ErrInvalidID = errors.New("id provided was not a valid UUID")
	// ErrForbidden occurs when a user tries to do something that is forbidden to
	// them according to our access control policies.
	ErrForbidden = errors.New("Attempted action is not allowed")
)


func List(ctx context.Context, db *sqlx.DB) ([]Product, error) {
    //list products

    const q = `SELECT
        p.product_id, p.name, p.cost, p.quantity,
        COALESCE(SUM(s.quantity),0) as sold,
        COALESCE(SUM(s.paid), 0) as revenue,
        p.date_updated, p.date_created
        FROM products AS p
        LEFT JOIN sales AS s on p.product_id = s.product_id
        GROUP BY p.product_id
        `;

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

    const q = `SELECT
        p.product_id, p.name, p.cost, p.quantity,
        COALESCE(SUM(s.quantity),0) as sold,
        COALESCE(SUM(s.paid), 0) as revenue,
        p.date_updated, p.date_created
        FROM products AS p
        LEFT JOIN sales AS s on p.product_id = s.product_id
        WHERE p.product_id = $1
        GROUP BY p.product_id
        `;

    if err := db.GetContext(ctx, &p, q, id); err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrNotFound
        }
        return nil, err
    }
    return &p, nil

}

func Create(ctx context.Context, db *sqlx.DB, user auth.Claims, np NewProduct, now time.Time) (*Product, error) {
    p := Product {
        ID:  uuid.New().String(),
        Name: np.Name,
        Cost: np.Cost,
        UserID: user.Subject,
        Quantity: np.Quantity,
        DateCreated: now.UTC(),
        DateUpdated: now.UTC(),
    }


    const q = "INSERT into products (product_id, name, cost, quantity, user_id, date_created, date_updated) VALUES ($1,$2,$3,$4,$5,$6,$7) ";

    if _, err := db.ExecContext(ctx, q, p.ID, p.Name, p.Cost, p.Quantity, p.UserID, p.DateCreated, p.DateUpdated); err != nil {
        return nil, errors.Wrapf(err, "insert product %v", np)
    }

    return &p, nil
}

// Update modifies data about a Product. It will error if the specified ID is
// invalid or does not reference an existing Product.
func Update(ctx context.Context, db *sqlx.DB, user auth.Claims, id string, update UpdateProduct, now time.Time) error {
	p, err := Retrieve(ctx, db, id)
	if err != nil {
		return err
	}

	// If you do not have the admin role ...
	// and you are not the owner of this product ...
	// then get outta here!
	if !user.HasRole(auth.RoleAdmin) && p.UserID != user.Subject {
		return ErrForbidden
	}
	if update.Name != nil {
		p.Name = *update.Name
	}
	if update.Cost != nil {
		p.Cost = *update.Cost
	}
	if update.Quantity != nil {
		p.Quantity = *update.Quantity
	}
	p.DateUpdated = now

	const q = `UPDATE products SET
		"name" = $2,
		"cost" = $3,
		"quantity" = $4,
		"date_updated" = $5
		WHERE product_id = $1`
	_, err = db.ExecContext(ctx, q, id,
		p.Name, p.Cost,
		p.Quantity, p.DateUpdated,
	)
	if err != nil {
		return errors.Wrap(err, "updating product")
	}

	return nil
}


// Delete product by id
func Delete(ctx context.Context, db *sqlx.DB, id string) error{
    if _,err := uuid.Parse(id); err != nil {
        return ErrInvalidID
    }

    const q = "delete from products where product_id = $1"

    if _, err := db.ExecContext(ctx, q, id); err != nil {
        return errors.Wrapf(err, "deleting product %s", id)
    }
    return nil
}
