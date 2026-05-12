package email

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

type MockEmailSender struct{}

func (m *MockEmailSender) Send(ctx context.Context, to string, subject string, body string) error {

	time.Sleep(500 * time.Millisecond)

	rand.Seed(time.Now().UnixNano())

	if rand.Float32() < 0.2 {
		return errors.New("provider temporary unavailable")
	}
	return nil
}
