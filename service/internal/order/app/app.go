package app

import (
	"AP2_assignment1/service/internal/order/repository"
	grpcserver "AP2_assignment1/service/internal/order/transport/grpc"
	httphandler "AP2_assignment1/service/internal/order/transport/http"
	"AP2_assignment1/service/internal/order/usecase"
	"database/sql"
	"log"
)

type OrderApp struct {
	Handler     *httphandler.OrderHandler
	UseCase     *usecase.OrderUsecase
	GRPCHandler *grpcserver.OrderGRPCHandler
}

func NewOrderApp(db *sql.DB, cfg *Config) *OrderApp {

	repository.InitRedis(cfg.RedisAddr)
	log.Println("Redis initialized!")

	repo := repository.NewOrderRepo(db)
	payClient := httphandler.NewPaymentGRPCClient(cfg.PaymentGRPCAddress)

	uc := usecase.NewOrderUsecase(repo, payClient)
	handler := httphandler.NewOrderHandler(uc)
	grpcHandler := grpcserver.NewOrderGRPCHandler(uc, repo)

	return &OrderApp{
		Handler:     handler,
		UseCase:     uc,
		GRPCHandler: grpcHandler,
	}
}
