package app

import (
	"AP2_assignment1/service/internal/payment/repository"
	"AP2_assignment1/service/internal/payment/transport/http"
	"AP2_assignment1/service/internal/payment/usecase"
	"database/sql"
)

type PaymentApp struct {
	Handler *http.PaymentHandler
}

func NewPaymentApp(db *sql.DB) *PaymentApp {
	repo := repository.NewPaymentRepo(db)
	uc := usecase.NewPaymentUsecase(repo)
	handler := http.NewPaymentHandler(uc)

	return &PaymentApp{
		Handler: handler,
	}
}
