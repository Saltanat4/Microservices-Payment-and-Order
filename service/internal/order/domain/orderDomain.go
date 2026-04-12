package domain

import "time"

type Order struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customer_id"`
	ItemName   string    `json:"item_name"`
	Amount     int64     `json:"amount"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type OrderUsecase interface {
	CreateOrder(ord *Order) error
}

type OrderRepository interface {
	Create(order *Order) error
	GetByID(id string) (*Order, error)
	UpdateStatus(id string, status string) error
	GetByAmountRange(min, max int64) ([]*Order, error)
	WatchStatusChange(id string, lastStatus string, done <-chan struct{}) (string, error)
}

type PaymentClient interface {
	Authorize(orderID string, amount int64) (string, error)
	Pay(orderID string, amount int64) (string, error)
}
