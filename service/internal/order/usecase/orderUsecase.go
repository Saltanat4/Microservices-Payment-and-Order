package usecase

import (
	"AP2_assignment1/service/internal/order/domain"
	"AP2_assignment1/service/internal/order/repository"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type OrderUsecase struct {
	repo      *repository.OrderRepo
	payClient domain.PaymentClient
}

func NewOrderUsecase(repo *repository.OrderRepo, payClient domain.PaymentClient) *OrderUsecase {
	return &OrderUsecase{repo: repo, payClient: payClient}
}

func (uc *OrderUsecase) CreateOrder(customerID, itemName string, amount int64, customer_email string) (*domain.Order, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	order := &domain.Order{
		ID:            uuid.New().String(),
		CustomerID:    customerID,
		CustomerEmail: customer_email,
		ItemName:      itemName,
		Amount:        amount,
		Status:        "Pending",
		CreatedAt:     time.Now(),
	}

	err := uc.repo.Create(order)
	if err != nil {
		return nil, err
	}

	payStatus, err := uc.payClient.Pay(order.ID, order.Amount, order.CustomerEmail)

	if err != nil || payStatus == "" {
		err := uc.repo.UpdateStatus(order.ID, "Failed")
		if err != nil {
			return nil, err
		}
		order.Status = "Failed"
		return order, nil
	}

	if payStatus == "Declined" {
		uc.repo.UpdateStatus(order.ID, "Failed")
		order.Status = "Failed"
		return order, nil
	}

	uc.repo.UpdateStatus(order.ID, "Paid")
	order.Status = "Paid"
	return order, nil
}

func (uc *OrderUsecase) GetOrder(id string) (*domain.Order, error) {
	ctx := context.Background()
	key := "order:" + id

	val, err := repository.RDB.Get(ctx, key).Result()
	if err == nil {
		var order domain.Order
		if err := json.Unmarshal([]byte(val), &order); err == nil {
			return &order, nil
		}
	}

	order, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(order)
	repository.RDB.Set(ctx, key, data, 5*time.Minute)

	return order, nil
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

	return uc.repo.GetByAmountRange(min, max)
}
