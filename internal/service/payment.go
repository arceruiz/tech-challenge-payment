package service

import (
	"context"
	"tech-challenge-payment/internal/canonical"
	"tech-challenge-payment/internal/config"
	"tech-challenge-payment/internal/integration/sqs_publisher"
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
	repo          repository.PaymentRepository
	publisher     sqs_publisher.Publisher
	statusToQueue map[canonical.PaymentStatus]string
}

func NewPaymentService() PaymentService {
	return &paymentService{
		repo:      repository.NewPaymentRepo(),
		publisher: sqs_publisher.NewSQS(),
		statusToQueue: map[canonical.PaymentStatus]string{
			canonical.PAYMENT_FAILED: config.Get().SQS.PaymentCancelledQueue,
			canonical.PAYMENT_PAYED:  config.Get().SQS.PaymentPayedQueue,
		},
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

	err = s.publisher.SendMessage(payment.OrderID, s.statusToQueue[status])
	if err != nil {
		return err
	}

	return nil
}

func (s *paymentService) GetAll(ctx context.Context) ([]canonical.Payment, error) {
	return s.repo.GetAll(ctx)
}
