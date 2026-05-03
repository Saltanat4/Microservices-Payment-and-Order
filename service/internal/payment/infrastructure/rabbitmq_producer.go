package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQProducer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQProducer(url string) (*RabbitMQProducer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	_ = ch.ExchangeDeclare("payment.dlx", "direct", true, false, false, false, nil)
	_, _ = ch.QueueDeclare("payment.dlq", true, false, false, false, nil)
	_ = ch.QueueBind("payment.dlq", "payment.dlq", "payment.dlx", false, nil)

	args := amqp.Table{
		"x-dead-letter-exchange":    "payment.dlx",
		"x-dead-letter-routing-key": "payment.dlq",
	}

	_, err = ch.QueueDeclare(
		"payment.completed",
		true,
		false,
		false,
		false,
		args,
	)

	return &RabbitMQProducer{conn: conn, channel: ch}, err
}

func (p *RabbitMQProducer) PublishPaymentEvent(orderID string, amount int64, status string, email string) error {
	body, _ := json.Marshal(map[string]interface{}{
		"order_id": orderID,
		"amount":   amount,
		"status":   status,
		"email":    email,
	})

	log.Printf("[Payment Service] Publishing event to RabbitMQ: OrderID=%s, Email=%s", orderID, email)

	return p.channel.PublishWithContext(context.Background(),
		"",
		"payment.completed",
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		})
}
