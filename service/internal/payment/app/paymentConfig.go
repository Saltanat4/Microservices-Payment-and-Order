package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort    string `env:"PAYMENT_HTTP_PORT"`
	GRPCPort    string `env:"PAYMENT_GRPC_PORT"`
	DBConn      string
	RabbitMQURL string
}

func NewConfig() *Config {
	_ = godotenv.Load(".env")

	host := strings.TrimSpace(os.Getenv("DB_HOST"))
	port := strings.TrimSpace(os.Getenv("DB_PORT_PAYMENT"))
	user := strings.TrimSpace(os.Getenv("DB_USER"))
	pass := strings.TrimSpace(os.Getenv("DB_PASSWORD"))
	name := strings.TrimSpace(os.Getenv("DB_NAME_PAYMENT"))

	host = "localhost"

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, name)

	httpPort := strings.TrimSpace(os.Getenv("PAYMENT_HTTP_PORT"))
	if httpPort == "" {
		httpPort = "8080"
	}

	grpcPort := strings.TrimSpace(os.Getenv("PAYMENT_GRPC_PORT"))
	if grpcPort == "" {
		grpcPort = "50051"
	}

	return &Config{
		HTTPPort:    httpPort,
		GRPCPort:    grpcPort,
		DBConn:      dsn,
		RabbitMQURL: rabbitURL,
	}
}
