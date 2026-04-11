package app

import (
	"AP2_assignment1/service/internal/order/repository"
	"AP2_assignment1/service/internal/order/transport/http"
	"AP2_assignment1/service/internal/order/usecase"
	"database/sql"
)

type OrderApp struct {
	Handler *http.OrderHandler
}

func NewOrderApp(db *sql.DB, cfg *Config) *OrderApp {
	repo := repository.NewOrderRepo(db)
	payClient := http.NewPaymentClient(cfg.PaymentService)

	uc := usecase.NewOrderUsecase(repo, payClient)

	handler := http.NewOrderHandler(uc)

	return &OrderApp{
		Handler: handler,
	}
}
