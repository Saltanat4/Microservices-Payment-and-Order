package app

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port   string
	DBConn string
}

func NewConfig() *Config {
	_ = godotenv.Load(".env")

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME_PAYMENT")

	host = "localhost"

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, name)

	return &Config{
		Port:   "8080",
		DBConn: dsn,
	}
}
