package usecase

import (
	"AP2_assignment1/service/internal/payment/domain"
	"AP2_assignment1/service/internal/payment/repository"
	"errors"

	"github.com/google/uuid"
)

type PaymentUsecase struct {
	repo     *repository.PaymentRepo
	producer domain.MessageProducer
}

func NewPaymentUsecase(repo *repository.PaymentRepo, producer domain.MessageProducer) *PaymentUsecase {
	return &PaymentUsecase{repo: repo, producer: producer}
}

func (uc *PaymentUsecase) Process(orderID string, amount int64, email string) (*domain.Payment, error) {
	status := "Authorized"
	if amount > 100000 {
		status = "Declined"
	}

	payment := &domain.Payment{
		OrderID:       orderID,
		Amount:        amount,
		Status:        status,
		CustomerEmail: email,
		TransactionID: uuid.New().String(),
	}

	_ = uc.repo.Save(payment)

	if payment.Status == "Authorized" && email != "" {
		_ = uc.producer.PublishPaymentEvent(payment.OrderID, payment.Amount, payment.Status, payment.CustomerEmail)
	}

	return payment, nil
}

func (uc *PaymentUsecase) GetStatus(orderID string) (*domain.Payment, error) {
	return uc.repo.GetByOrderID(orderID)
}

func (uc *PaymentUsecase) ListPayments(status string) ([]*domain.Payment, error) {
	if status == " " {
		return nil, errors.New("invalid")
	}
	return uc.repo.ListByStatus(status)
}
