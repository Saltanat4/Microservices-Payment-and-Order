package http

import (
	"AP2_assignment1/service/internal/order/domain"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type PaymentClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewPaymentClient(url string) domain.PaymentClient {
	return &PaymentClient{
		baseURL:    url,
		httpClient: &http.Client{Timeout: 2 * time.Second},
	}
}
func (c *PaymentClient) Authorize(orderID string, amount int64) (string, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"order_id": orderID,
		"amount":   amount,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := c.httpClient.Post(c.baseURL+"/payments", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("payment service unavailable: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		TransactionID string `json:"transaction_id"`
		Status        string `json:"status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode payment response: %v", err)
	}

	if result.Status == "Declined" {
		return "", fmt.Errorf("payment declined")
	}

	return result.TransactionID, nil
}

func (c *PaymentClient) Pay(orderID string, amount int64) (string, error) {

	jsonData := fmt.Sprintf(`{"order_id":"%s", "amount":%d}`, orderID, amount)
	resp, err := c.httpClient.Post(c.baseURL+"/payments", "application/json", strings.NewReader(jsonData))

	if err != nil {
		return "", errors.New("payment service unavailable")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "Declined", nil
	}

	return "Completed", nil
}
