package http

import (
	"AP2_assignment1/service/internal/order/domain"
	"context"
	"fmt"
	"log"
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
		log.Printf("WARNING: could not connect to payment gRPC at %s: %v", address, err)
	}

	return &PaymentGRPCClient{
		conn:   conn,
		client: pb.NewPaymentServiceClient(conn),
	}
}

func (c *PaymentGRPCClient) Authorize(orderID string, amount int64) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := c.client.ProcessPayment(ctx, &pb.PaymentRequest{
		OrderId: orderID,
		Amount:  amount,
	})
	if err != nil {
		return "", fmt.Errorf("payment service unavailable: %v", err)
	}

	if resp.Status == "Declined" {
		return "", fmt.Errorf("payment declined")
	}

	return resp.TransactionId, nil
}

func (c *PaymentGRPCClient) Pay(orderID string, amount int64) (string, error) {
	tid, err := c.Authorize(orderID, amount)
	if err != nil {
		return "Declined", nil
	}
	return tid, nil
}
