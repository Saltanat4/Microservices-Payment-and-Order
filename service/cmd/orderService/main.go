package main

import (
	"AP2_assignment1/service/internal/order/app"
	"database/sql"
	"log"
	"net"

	pb "github.com/Saltanat4/gen-repo/order"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	cfg := app.NewConfig()

	db, err := sql.Open("postgres", cfg.DBConn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	orderApp := app.NewOrderApp(db, cfg)

	go func() {
		lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
		if err != nil {
			log.Fatalf("failed to listen on gRPC port %s: %v", cfg.GRPCPort, err)
		}

		s := grpc.NewServer()
		pb.RegisterOrderServiceServer(s, orderApp.GRPCHandler)

		log.Printf("Order gRPC Server (streaming) starting on :%s", cfg.GRPCPort)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	r := gin.Default()
	r.POST("/orders", orderApp.Handler.CreateOrder)
	r.GET("/orders/:id", orderApp.Handler.GetOrder)
	r.PUT("/orders/:id/cancel", orderApp.Handler.CancelOrder)
	r.GET("/orders", orderApp.Handler.GetOrders)

	log.Printf("Order HTTP Server starting on :%s", cfg.HTTPPort)
	if err := r.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}
