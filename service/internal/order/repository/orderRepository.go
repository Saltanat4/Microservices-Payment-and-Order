package repository

import (
	"AP2_assignment1/service/internal/order/domain"
	"context"
	"database/sql"
	"time"
)

type OrderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) Create(o *domain.Order) error {
	query := `INSERT INTO orders (id, customer_id, item_name, amount, status, customer_email, created_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(query, o.ID, o.CustomerID, o.ItemName, o.Amount, o.CustomerEmail, o.Status, o.CreatedAt)
	return err
}

func (r *OrderRepo) GetByID(id string) (*domain.Order, error) {
	o := &domain.Order{}
	query := `SELECT id, customer_id, item_name, amount, status, customer_email,created_at
              FROM orders WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&o.ID,
		&o.CustomerID,
		&o.ItemName,
		&o.Amount,
		&o.Status,
		&o.CustomerEmail,
		&o.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return o, nil
}

func (r *OrderRepo) UpdateStatus(id string, status string) error {

	_, err := r.db.Exec("UPDATE orders SET status=$1 WHERE id=$2", status, id)
	if err != nil {
		return err
	}

	ctx := context.Background()
	RDB.Del(ctx, "order:"+id)

	return nil
}
func (r *OrderRepo) GetByAmountRange(min, max int64) ([]*domain.Order, error) {
	query := `SELECT id, customer_id, item_name, amount, status, customer_email, created_at
              FROM orders
              WHERE amount >= $1 AND amount <= $2`

	rows, err := r.db.Query(query, min, max)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		o := &domain.Order{}
		err := rows.Scan(&o.ID, &o.CustomerID, &o.ItemName, &o.Amount, &o.Status, &o.CustomerEmail, &o.CreatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *OrderRepo) WatchStatusChange(id string, lastStatus string, done <-chan struct{}) (string, error) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return "", nil
		case <-ticker.C:
			var currentStatus string
			query := `SELECT status FROM orders WHERE id = $1`
			err := r.db.QueryRow(query, id).Scan(&currentStatus)
			if err != nil {
				return "", err
			}
			if currentStatus != lastStatus {
				return currentStatus, nil
			}
		}
	}
}
