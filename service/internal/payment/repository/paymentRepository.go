package repository

import (
	"AP2_assignment1/service/internal/payment/domain"
	"database/sql"
)

type PaymentRepo struct {
	db *sql.DB
}

func NewPaymentRepo(db *sql.DB) *PaymentRepo {
	return &PaymentRepo{db: db}
}

func (r *PaymentRepo) Save(p *domain.Payment) error {
	query := `INSERT INTO payments (order_id, amount, status, transaction_id) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, p.OrderID, p.Amount, p.Status, p.TransactionID)
	return err
}

func (r *PaymentRepo) GetByOrderID(orderID string) (*domain.Payment, error) {
	p := &domain.Payment{}
	query := `SELECT order_id, amount, status, transaction_id FROM payments WHERE order_id = $1`
	err := r.db.QueryRow(query, orderID).Scan(&p.OrderID, &p.Amount, &p.Status, &p.TransactionID)
	return p, err
}
