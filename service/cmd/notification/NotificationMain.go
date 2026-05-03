package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	processedOrders = make(map[string]bool)
	mutex           = &sync.Mutex{}
)

func main() {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	//_, err = ch.QueueDeclare(
	//"payment.completed",
	//true,
	//false,
	//false,
	//false,
	//nil,
	//)

	msgs, err := ch.Consume(
		"payment.completed",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Println("Waiting for payment messages...")

	for d := range msgs {
		var event struct {
			OrderID string `json:"order_id"`
			Email   string `json:"email"`
			Amount  int64  `json:"amount"`
		}

		_ = json.Unmarshal(d.Body, &event)
		if event.Amount == 101 {
			log.Println("Permanent error for Order #%s: Invalid amount (%d)", event.OrderID, event.Amount)
			log.Println("Moving message...")

			d.Nack(false, false)
			continue
		}

		if err := json.Unmarshal(d.Body, &event); err != nil {
			log.Printf("Error decoding message: %v", err)
			d.Ack(false)
			continue
		}

		mutex.Lock()
		if processedOrders[event.OrderID] {
			log.Printf("Duplicate message ignored for Order #%s", event.OrderID)
			mutex.Unlock()
			d.Ack(false)
			continue
		}

		processedOrders[event.OrderID] = true
		mutex.Unlock()

		log.Printf("[Notification] Sent email to %s for Order #%s. Amount: %d",
			event.Email, event.OrderID, event.Amount)

		log.Printf("[Notification] Done! Sending ACK.")
		d.Ack(false)
	}
}
