package order

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Repository interface {
	Close()
	PutOrder(ctx context.Context, o Order) error
	GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &postgresRepository{db}, nil
}

func (r *postgresRepository) Close() {
	r.db.Close()
}

func (r *postgresRepository) PutOrder(ctx context.Context, o Order) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	// Insert order in orders table
	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO orders(id, created_at, account_id, total_price) VALUES($1,$2,$3,$4)",
		o.ID,
		o.CreatedAt,
		o.AccountID,
		o.TotalPrice,
	)
	if err != nil {
		return
	}

	// Insert ordered products in order_products table
	stmt, _ := tx.PrepareContext(ctx, pq.CopyIn("order_products", "order_id", "product_id", "quantity"))
	for _, p := range o.Products {
		_, err = tx.ExecContext(ctx, o.ID, p.ID, p.Quantity)
		if err != nil {
			return
		}
	}
	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return
	}
	stmt.Close()

	return
}

func (r *postgresRepository) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT
		o.id,
		o.created_at,
		o.account_id,
		o.total_price::money::numeric::float8,
      	op.product_id,
      	op.quantity
		FROM orders o JOIN order_products op ON (o.id = op.order_id)
		WHERE o.account_id = $1
		ORDER BY o.id`,
		accountID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ordersByID := make(map[string]*Order)
	orderIDs := []string{}

	for rows.Next() {
		var (
			orderID    string
			createdAt  time.Time
			accID      string
			totalPrice float64
			productID  string
			quantity   uint32
		)

		if err := rows.Scan(
			&orderID,
			&createdAt,
			&accID,
			&totalPrice,
			&productID,
			&quantity,
		); err != nil {
			return nil, err
		}

		order, exists := ordersByID[orderID]
		if !exists {
			order = &Order{
				ID:         orderID,
				AccountID:  accID,
				CreatedAt:  createdAt,
				TotalPrice: totalPrice,
			}
			ordersByID[orderID] = order
			orderIDs = append(orderIDs, orderID)
		}

		order.Products = append(order.Products, OrderedProduct{
			ID:       productID,
			Quantity: quantity,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// preserve the order request
	orders := make([]Order, 0, len(orderIDs))
	for _, id := range orderIDs {
		orders = append(orders, *ordersByID[id])
	}

	return orders, nil
}
