package domain

type Payment struct {
	ID            int    `json:"-"`
	OrderID       string `json:"order_id"`
	Amount        int64  `json:"amount"`
	Status        string `json:"status"`
	CustomerEmail string `json:"customer_email"`
	TransactionID string `json:"transaction_id"`
}

type MessageProducer interface {
	PublishPaymentEvent(orderID string, amount int64, status string, email string) error
}
