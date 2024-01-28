package service

import (
	"context"
	"tech-challenge-payment/internal/canonical"
	"tech-challenge-payment/internal/integration/order"
	"tech-challenge-payment/internal/repository"
	"time"
)

const (
	ORDER_PAYED     = "PAYED"
	ORDER_CANCELLED = "CANCELLED"
)

type PaymentService interface {
	GetByID(context.Context, string) (*canonical.Payment, error)
	Callback(ctx context.Context, paymentId string, status canonical.PaymentStatus) error
	Create(ctx context.Context, payment canonical.Payment) (*canonical.Payment, error)
	GetAll(ctx context.Context) ([]canonical.Payment, error)
}

type paymentService struct {
	repo         repository.PaymentRepository
	orderService order.OrderService
}

func NewPaymentService(repo repository.PaymentRepository, orderService order.OrderService) PaymentService {
	return &paymentService{
		repo:         repo,
		orderService: orderService,
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

	orderStatus := ""
	switch status {
	case canonical.PAYMENT_PAYED:
		orderStatus = ORDER_PAYED
	case canonical.PAYMENT_FAILED:
		orderStatus = ORDER_CANCELLED
	default:
		return nil
	}

	err = s.orderService.UpdateStatus(payment.OrderID, orderStatus)
	if err != nil {
		return err
	}

	return nil
}

func (s *paymentService) GetAll(ctx context.Context) ([]canonical.Payment, error) {
	return s.repo.GetAll(ctx)
}
