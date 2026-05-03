package app

import (
	"AP2_assignment1/service/internal/payment/infrastructure"
	"AP2_assignment1/service/internal/payment/repository"
	httphandler "AP2_assignment1/service/internal/payment/transport/http"
	"AP2_assignment1/service/internal/payment/usecase"
	"database/sql"
	"log"
)

type PaymentApp struct {
	Handler *httphandler.PaymentHandler
	UseCase *usecase.PaymentUsecase
}

func NewPaymentApp(db *sql.DB, cfg *Config) *PaymentApp {
	repo := repository.NewPaymentRepo(db)

	producer, err := infrastructure.NewRabbitMQProducer(cfg.RabbitMQURL)
	if err != nil {
		log.Printf("Warning: RabbitMQ not connected: %v", err)
	}

	uc := usecase.NewPaymentUsecase(repo, producer)
	handler := httphandler.NewPaymentHandler(uc)

	return &PaymentApp{
		Handler: handler,
		UseCase: uc,
	}
}
