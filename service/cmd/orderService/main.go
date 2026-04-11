package main

import (
	"AP2_assignment1/service/internal/order/app"
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	cfg := app.NewConfig()

	db, err := sql.Open("postgres", cfg.DBConn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("ERROR ORDER_DB: %v", err)
	} else {
		log.Println("SUCCESS!")
	}

	orderApp := app.NewOrderApp(db, cfg)

	r := gin.Default()
	r.POST("/orders", orderApp.Handler.CreateOrder)
	r.GET("/orders", orderApp.Handler.GetOrders)
	r.GET("/orders/:id", orderApp.Handler.GetOrder)
	r.PATCH("/orders/:id/cancel", orderApp.Handler.CancelOrder)

	r.Run(":" + cfg.Port)
}

//8081-order 8080-payment
