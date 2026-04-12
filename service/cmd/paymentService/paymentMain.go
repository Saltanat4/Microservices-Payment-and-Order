package main

import (
	"AP2_assignment1/service/internal/payment/app"
	grpchandler "AP2_assignment1/service/internal/payment/transport/grpc"
	"database/sql"
	"log"
	"net"

	pb "github.com/Saltanat4/gen-repo/payment"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	cfg := app.NewConfig()

	db, err := sql.Open("postgres", cfg.DBConn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	paymentApp := app.NewPaymentApp(db)

	go func() {
		lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
		if err != nil {
			log.Fatalf("failed to listen on gRPC port %s: %v", cfg.GRPCPort, err)
		}

		s := grpc.NewServer()
		grpcHandler := grpchandler.NewPaymentGRPCHandler(paymentApp.UseCase)
		pb.RegisterPaymentServiceServer(s, grpcHandler)

		log.Printf("Payment gRPC Server starting on :%s", cfg.GRPCPort)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	r := gin.Default()
	r.POST("/payments", paymentApp.Handler.ProcessPayment)
	r.GET("/payments/:order_id", paymentApp.Handler.GetPaymentStatus)

	log.Printf("Payment HTTP Server starting on :%s", cfg.HTTPPort)
	if err := r.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}
