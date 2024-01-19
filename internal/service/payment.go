package service

import (
	"context"
	"tech-challenge-payment/internal/canonical"
	"tech-challenge-payment/internal/repository"
	"time"
)

type PaymentService interface {
	GetByID(context.Context, string) (*canonical.Payment, error)
	Callback(ctx context.Context, paymentId string, status canonical.PaymentStatus) error
	Create(ctx context.Context, payment canonical.Payment) (*canonical.Payment, error)
}

type paymentService struct {
	repo repository.PaymentRepository
}

func NewPaymentService(repo repository.PaymentRepository) PaymentService {
	return &paymentService{
		repo: repo,
	}
}

func (s *paymentService) Create(ctx context.Context, payment canonical.Payment) (*canonical.Payment, error) {
	payment.Status = canonical.PAYMENT_CREATED
	payment.ID = canonical.NewUUID()
	payment.CreatedAt = time.Now()
	payment, err := s.repo.Create(ctx, payment)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (s *paymentService) GetByID(ctx context.Context, id string) (*canonical.Payment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *paymentService) Callback(ctx context.Context, paymentId string, status canonical.PaymentStatus) error {
	payment, err := s.repo.GetByID(ctx, paymentId)
	if err != nil {
		return err
	}
	if payment == nil {
		return canonical.ErrorNotFound
	}

	payment.UpdatedAt = time.Now()
	payment.Status = status
	err = s.repo.Update(ctx, paymentId, *payment)
	if err != nil {
		return err
	}
	return nil
}
