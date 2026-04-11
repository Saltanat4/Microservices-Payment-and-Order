package main

import (
	"AP2_assignment1/service/internal/payment/app"
	"database/sql"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	cfg := app.NewConfig()
	db, err := sql.Open("postgres", cfg.DBConn)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Database not available: %v", err)
	}
	paymentApp := app.NewPaymentApp(db)

	r := gin.Default()

	time.Sleep(5 * time.Second)
	r.POST("/payments", paymentApp.Handler.ProcessPayment)
	r.GET("/payments/:order_id", paymentApp.Handler.GetPaymentStatus)

	r.Run(":" + cfg.Port)
	log.Printf("Payment Service %s", cfg.Port)
}
