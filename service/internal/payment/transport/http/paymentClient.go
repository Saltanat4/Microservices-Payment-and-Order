package http

import (
	"AP2_assignment1/service/internal/order/domain"
	"context"
	"fmt"
	"time"

	pb "github.com/Saltanat4/gen-repo/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PaymentGRPCClient struct {
	client pb.PaymentServiceClient
	conn   *grpc.ClientConn
}

func NewPaymentGRPCClient(address string) domain.PaymentClient {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &PaymentGRPCClient{
			conn:   nil,
			client: nil,
		}
	}

	return &PaymentGRPCClient{
		conn:   conn,
		client: pb.NewPaymentServiceClient(conn),
	}
}

func (c *PaymentGRPCClient) Authorize(orderID string, amount int64, email string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := c.client.ProcessPayment(ctx, &pb.PaymentRequest{
		OrderId:       orderID,
		Amount:        amount,
		CustomerEmail: email,
	})

	if err != nil {
		return "", err
	}

	return resp.TransactionId, nil
}

func (c *PaymentGRPCClient) Pay(orderID string, amount int64, email string) (string, error) {
	if c.client == nil {
		return "", fmt.Errorf("payment service unavailable")
	}

	return c.Authorize(orderID, amount, email)
}
