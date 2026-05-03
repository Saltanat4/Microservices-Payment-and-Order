package grpc

import (
	"AP2_assignment1/service/internal/payment/usecase"
	"context"
	"time"

	pb "github.com/Saltanat4/gen-repo/payment"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PaymentGRPCHandler struct {
	pb.UnimplementedPaymentServiceServer
	uc *usecase.PaymentUsecase
}

func NewPaymentGRPCHandler(uc *usecase.PaymentUsecase) *PaymentGRPCHandler {
	return &PaymentGRPCHandler{uc: uc}
}

func (h *PaymentGRPCHandler) ListPayments(
	ctx context.Context,
	req *pb.ListPaymentsRequest,
) (*pb.ListPaymentsResponse, error) {

	payments, err := h.uc.ListPayments(req.Status)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var result []*pb.PaymentResponse

	for _, p := range payments {
		result = append(result, &pb.PaymentResponse{
			OrderId:       p.OrderID,
			Amount:        p.Amount,
			Status:        p.Status,
			TransactionId: p.TransactionID,
			ProcessedAt:   timestamppb.New(time.Now()),
		})
	}

	return &pb.ListPaymentsResponse{
		Payments: result,
	}, nil
}

func (h *PaymentGRPCHandler) ProcessPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	result, err := h.uc.Process(req.OrderId, req.Amount, req.CustomerEmail)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.PaymentResponse{
		TransactionId: result.TransactionID,
		OrderId:       result.OrderID,
		Amount:        result.Amount,
		Status:        result.Status,
		ProcessedAt:   timestamppb.New(time.Now()),
	}, nil
}
