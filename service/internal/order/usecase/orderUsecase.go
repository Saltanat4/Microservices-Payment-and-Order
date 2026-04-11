package usecase

import (
	"AP2_assignment1/service/internal/order/domain"
	"AP2_assignment1/service/internal/order/repository"
	"errors"
	"time"

	"github.com/google/uuid"
)

type PaymentProvider interface {
	Pay(orderID string, amount int64) (string, error)
}

type OrderUsecase struct {
	repo      *repository.OrderRepo
	payClient domain.PaymentClient
}

func NewOrderUsecase(repo *repository.OrderRepo, payClient domain.PaymentClient) *OrderUsecase {
	return &OrderUsecase{repo: repo, payClient: payClient}
}

func (uc *OrderUsecase) CreateOrder(customerID, itemName string, amount int64) (*domain.Order, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}
	order := &domain.Order{
		ID:         uuid.New().String(),
		CustomerID: customerID,
		ItemName:   itemName,
		Amount:     amount,
		Status:     "Pending",
		CreatedAt:  time.Now(),
	}

	err := uc.repo.Create(order)
	if err != nil {
		return nil, err
	}

	payStatus, err := uc.payClient.Pay(order.ID, order.Amount)

	if err != nil {
		return nil, errors.New("payment_service_unavailable")
	}

	if payStatus != "Completed" {
		uc.repo.UpdateStatus(order.ID, "Failed")
		order.Status = "Failed"
		return order, nil
	}

	uc.repo.UpdateStatus(order.ID, "Paid")
	order.Status = "Paid"

	return order, nil
}

func (uc *OrderUsecase) GetOrder(id string) (*domain.Order, error) {

	return uc.repo.GetByID(id)
}

func (uc *OrderUsecase) CancelOrder(id string) error {
	order, err := uc.repo.GetByID(id)
	if err != nil {
		return err
	}

	if order.Status != "Pending" {
		return errors.New("prohibited: only Pending orders can be cancelled")
	}

	return uc.repo.UpdateStatus(id, "Cancelled")
}

func (uc *OrderUsecase) GetOrdersByRange(min, max int64) ([]*domain.Order, error) {
	if min < 0 || min >= max {
		return nil, errors.New("bad_request: invalid range")
	}

	orders, err := uc.repo.GetByAmountRange(min, max)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
