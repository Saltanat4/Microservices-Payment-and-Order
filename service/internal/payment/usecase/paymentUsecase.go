package usecase

import (
	"AP2_assignment1/service/internal/payment/domain"
	"AP2_assignment1/service/internal/payment/repository"

	"github.com/google/uuid"
)

type PaymentUsecase struct {
	repo *repository.PaymentRepo
}

func NewPaymentUsecase(repo *repository.PaymentRepo) *PaymentUsecase {
	return &PaymentUsecase{repo: repo}
}

func (uc *PaymentUsecase) Process(orderID string, amount int64) (*domain.Payment, error) {
	status := "Authorized"
	if amount > 100000 {
		status = "Declined"
	}
	payment := &domain.Payment{
		OrderID:       orderID,
		Amount:        amount,
		Status:        status,
		TransactionID: uuid.New().String(),
	}

	err := uc.repo.Save(payment)
	return payment, err
}

func (uc *PaymentUsecase) GetStatus(orderID string) (*domain.Payment, error) {
	return uc.repo.GetByOrderID(orderID)
}
