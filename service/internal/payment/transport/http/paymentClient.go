package http

import (
	"net/http"
	"time"
)

type PaymentClient struct {
	httpClient *http.Client
}

func NewPaymentClient() *PaymentClient {
	return &PaymentClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}
