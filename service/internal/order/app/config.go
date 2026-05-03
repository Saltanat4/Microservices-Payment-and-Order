package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort           string `env:"ORDER_HTTP_PORT"`
	GRPCPort           string `env:"ORDER_GRPC_PORT"`
	DBConn             string
	PaymentGRPCAddress string `env:"ORDER_PAYMENT_GRPC_ADDRESS"`
	RabbitMQURL        string
}

func NewConfig() *Config {
	_ = godotenv.Load(".env")

	user := strings.TrimSpace(os.Getenv("DB_USER"))
	pass := strings.TrimSpace(os.Getenv("DB_PASSWORD"))
	host := strings.TrimSpace(os.Getenv("DB_HOST"))
	port := strings.TrimSpace(os.Getenv("DB_PORT_ORDER"))
	name := strings.TrimSpace(os.Getenv("DB_NAME_ORDER"))

	host = "localhost"

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, name)

	httpPort := strings.TrimSpace(os.Getenv("ORDER_HTTP_PORT"))
	if httpPort == "" {
		httpPort = "8081"
	}

	grpcPort := strings.TrimSpace(os.Getenv("ORDER_GRPC_PORT"))
	if grpcPort == "" {
		grpcPort = "50052"
	}

	paymentAddr := strings.TrimSpace(os.Getenv("PAYMENT_GRPC_ADDRESS"))
	if paymentAddr == "" {
		paymentAddr = "localhost:50051"
	}

	fmt.Printf("DEBUG ORDER DSN: '%s'\n", dsn)
	fmt.Printf("DEBUG PAYMENT_GRPC_ADDRESS: '%s'\n", paymentAddr)

	return &Config{
		HTTPPort:           httpPort,
		GRPCPort:           grpcPort,
		DBConn:             dsn,
		PaymentGRPCAddress: paymentAddr,
		RabbitMQURL:        rabbitURL,
	}
}
