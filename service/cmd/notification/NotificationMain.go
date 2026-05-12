package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"AP2_assignment1/service/internal/notification/app"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}

	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
	ctx := context.Background()

	notificationApp := app.NewNotificationApp()
	emailSender := notificationApp.EmailSender

	rabbitURL := strings.TrimSpace(os.Getenv("RABBITMQ_URL"))
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@127.0.0.1:5672/"
	}

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("RabbitMQ connection failed: %v", err)
	}
	defer conn.Close()

	ch, _ := conn.Channel()
	msgs, _ := ch.Consume("payment.completed", "", false, false, false, false, nil)

	log.Println("Notification Worker is running...")

	for d := range msgs {
		var event struct {
			OrderID string  `json:"order_id"`
			Email   string  `json:"email"`
			Amount  float64 `json:"amount"`
		}

		if err := json.Unmarshal(d.Body, &event); err != nil {
			d.Ack(false)
			continue
		}

		if event.Amount == 101 {
			log.Printf("Suspicious amount (101) for Order %s. Sending to DLQ.", event.OrderID)
			d.Nack(false, false)
			continue
		}

		key := "notified:" + event.OrderID
		isNew, _ := rdb.SetNX(ctx, key, "true", 24*time.Hour).Result()
		if !isNew {
			log.Printf("Duplicate order #%s ignored", event.OrderID)
			d.Ack(false)
			continue
		}

		maxRetries := 3
		retryDelay := 2 * time.Second
		var sendErr error

		for i := 0; i < maxRetries; i++ {
			sendErr = emailSender.Send(ctx, event.Email, "Order Confirmed", "Success!")
			if sendErr == nil {
				break
			}

			log.Printf("Attempt %d failed for order %s. Retrying in %v...", i+1, event.OrderID, retryDelay)
			time.Sleep(retryDelay)
			retryDelay *= 2
		}

		if sendErr != nil {
			log.Printf("Critical: Failed to notify order %s. Sending to DLQ.", event.OrderID)
			rdb.Del(ctx, key)
			d.Nack(false, false)
		} else {
			fmt.Printf("[Notification] Sent email to %s for Order #%s. Amount: $%.2f\n", event.Email, event.OrderID, event.Amount)
			d.Ack(false)
		}
	}
}
