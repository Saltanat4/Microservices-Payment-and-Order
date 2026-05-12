package main

import (
	"AP2_assignment1/service/internal/order/app"
	"AP2_assignment1/service/internal/order/transport/http"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	pb "github.com/Saltanat4/gen-repo/order"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

func main() {
	cfg := app.NewConfig()

	db, err := sql.Open("postgres", cfg.DBConn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}(db)

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	orderApp := app.NewOrderApp(db, cfg)

	go func() {
		lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
		if err != nil {
			log.Fatalf("failed to listen on gRPC port %s: %v", cfg.GRPCPort, err)
		}

		s := grpc.NewServer()
		pb.RegisterOrderServiceServer(s, orderApp.GRPCHandler)

		log.Printf("Order gRPC Server starting on :%s", cfg.GRPCPort)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	fmt.Println("ORDER SERVICE STARTED")

	r := gin.Default()

	limit := 3
	window := time.Minute
	r.POST("/orders", http.RateLimiterMiddleware(rdb, limit, window), orderApp.Handler.CreateOrder)
	r.GET("/orders/:id", orderApp.Handler.GetOrder)
	r.PATCH("/orders/:id/cancel", orderApp.Handler.CancelOrder)
	r.GET("/orders", http.RateLimiterMiddleware(rdb, limit, window), orderApp.Handler.GetOrders)
	log.Printf("Order HTTP Server starting on :%s", cfg.HTTPPort)
	if err := r.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}
