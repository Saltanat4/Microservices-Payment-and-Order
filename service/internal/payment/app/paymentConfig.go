package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort string
	GRPCPort string
	DBConn   string
}

func NewConfig() *Config {
	_ = godotenv.Load(".env")

	host := strings.TrimSpace(os.Getenv("DB_HOST"))
	port := strings.TrimSpace(os.Getenv("DB_PORT"))
	user := strings.TrimSpace(os.Getenv("DB_USER"))
	pass := strings.TrimSpace(os.Getenv("DB_PASSWORD"))
	name := strings.TrimSpace(os.Getenv("DB_NAME_PAYMENT"))

	host = "localhost"

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
		HTTPPort: httpPort,
		GRPCPort: grpcPort,
		DBConn:   dsn,
	}
}
