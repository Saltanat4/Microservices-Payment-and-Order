package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	pb "github.com/Saltanat4/gen-repo/order"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: stream_client <order_id>")
	}
	orderID := os.Args[1]

	_ = godotenv.Load(".env")

	orderGRPCAddr := strings.TrimSpace(os.Getenv("ORDER_GRPC_ADDRESS"))
	if orderGRPCAddr == "" {
		orderGRPCAddr = "localhost:50052"
	}

	conn, err := grpc.Dial(orderGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to Order gRPC server at %s: %v", orderGRPCAddr, err)
	}
	defer conn.Close()

	client := pb.NewOrderServiceClient(conn)

	stream, err := client.SubscribeToOrderUpdates(context.Background(), &pb.OrderRequest{
		OrderId: orderID,
	})
	if err != nil {
		log.Fatalf("failed to subscribe: %v", err)
	}

	fmt.Printf("Subscribed to order %s updates...\n", orderID)

	for {
		update, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("Stream closed by server (order reached terminal state).")
			break
		}
		if err != nil {
			log.Fatalf("error receiving update: %v", err)
		}

		fmt.Printf("[%s] Order %s -> Status: %s\n",
			update.UpdatedAt.AsTime().Format("15:04:05"),
			update.OrderId,
			update.Status,
		)
	}
}
