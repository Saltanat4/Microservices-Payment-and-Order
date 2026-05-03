package grpc

import (
	"AP2_assignment1/service/internal/order/repository"
	"log"
	"time"

	pb "github.com/Saltanat4/gen-repo/order"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderGRPCHandler struct {
	pb.UnimplementedOrderServiceServer
	repo *repository.OrderRepo
}

func NewOrderGRPCHandler(repo *repository.OrderRepo) *OrderGRPCHandler {
	return &OrderGRPCHandler{repo: repo}
}

func (h *OrderGRPCHandler) SubscribeToOrderUpdates(req *pb.OrderRequest, stream pb.OrderService_SubscribeToOrderUpdatesServer) error {
	orderID := req.OrderId

	order, err := h.repo.GetByID(orderID)
	if err != nil {
		return status.Errorf(codes.NotFound, "order not found: %v", err)
	}

	err = stream.Send(&pb.OrderStatusUpdate{
		OrderId:   order.ID,
		Status:    order.Status,
		UpdatedAt: timestamppb.New(time.Now()),
	})
	if err != nil {
		return status.Errorf(codes.Internal, "failed to send initial update: %v", err)
	}

	log.Printf("Client subscribed to order %s updates (current: %s)", orderID, order.Status)

	lastStatus := order.Status

	terminalStatuses := map[string]bool{
		"Paid":      true,
		"Failed":    true,
		"Cancelled": true,
	}

	if terminalStatuses[lastStatus] {
		log.Printf("Order %s is already in terminal state: %s. Closing stream.", orderID, lastStatus)
		return nil
	}

	for {
		done := stream.Context().Done()

		newStatus, err := h.repo.WatchStatusChange(orderID, lastStatus, done)
		if err != nil {
			if stream.Context().Err() != nil {
				log.Printf("Client disconnected from order %s stream", orderID)
				return nil
			}
			return status.Errorf(codes.Internal, "error watching status: %v", err)
		}

		if newStatus == "" {
			log.Printf("Client disconnected from order %s stream", orderID)
			return nil
		}

		log.Printf("Order %s status changed: %s -> %s", orderID, lastStatus, newStatus)

		err = stream.Send(&pb.OrderStatusUpdate{
			OrderId:   orderID,
			Status:    newStatus,
			UpdatedAt: timestamppb.New(time.Now()),
		})
		if err != nil {
			log.Printf("Failed to send update for order %s: %v", orderID, err)
			return err
		}

		lastStatus = newStatus

		if terminalStatuses[newStatus] {
			log.Printf("Order %s reached terminal state: %s. Closing stream.", orderID, newStatus)
			return nil
		}
	}
}
