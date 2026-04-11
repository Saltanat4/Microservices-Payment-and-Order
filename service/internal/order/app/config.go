package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DBConn         string
	PaymentService string
}

func NewConfig() *Config {
	_ = godotenv.Load(".env")

	user := strings.TrimSpace(os.Getenv("DB_USER"))
	pass := strings.TrimSpace(os.Getenv("DB_PASSWORD"))
	host := strings.TrimSpace(os.Getenv("DB_HOST"))
	port := strings.TrimSpace(os.Getenv("DB_PORT"))
	name := strings.TrimSpace(os.Getenv("DB_NAME_ORDER"))
	host = "localhost"
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, name)

	fmt.Printf("DEBUG ORDER DSN: '%s'\n", dsn)

	fmt.Println(host + " 1 " + port + " 2 " + user + " 3 " + pass + " 4 " + name)
	return &Config{
		Port:           "8081",
		PaymentService: "http://localhost:8080",
		DBConn:         dsn,
	}
}
