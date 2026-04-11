package domain

type Payment struct {
	ID            int    `json:"-"`
	OrderID       string `json:"order_id"`
	Amount        int64  `json:"amount"`
	Status        string `json:"status"`
	TransactionID string `json:"transaction_id"`
}
