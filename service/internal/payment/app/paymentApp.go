package app

import (
	"AP2_assignment1/service/internal/payment/repository"
	httphandler "AP2_assignment1/service/internal/payment/transport/http"
	"AP2_assignment1/service/internal/payment/usecase"
	"database/sql"
)

type PaymentApp struct {
	Handler *httphandler.PaymentHandler
	UseCase *usecase.PaymentUsecase
}

func NewPaymentApp(db *sql.DB) *PaymentApp {
	repo := repository.NewPaymentRepo(db)
	uc := usecase.NewPaymentUsecase(repo)
	handler := httphandler.NewPaymentHandler(uc)

	return &PaymentApp{
		Handler: handler,
		UseCase: uc,
	}
}
