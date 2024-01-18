package repository

import (
	"context"
	"tech-challenge-payment/internal/canonical"

	"github.com/stretchr/testify/mock"
)

type PaymentRepositoryMock struct {
	mock.Mock
}

func (m *PaymentRepositoryMock) Create(ctx context.Context, payment canonical.Payment) (canonical.Payment, error) {
	args := m.Called(ctx, payment)
	return args.Get(0).(canonical.Payment), args.Error(1)
}

func (m *PaymentRepositoryMock) Update(ctx context.Context, paymentId string, payment canonical.Payment) error {
	args := m.Called(ctx, paymentId, payment)
	return args.Error(0)
}

func (m *PaymentRepositoryMock) GetByID(ctx context.Context, paymentId string) (*canonical.Payment, error) {
	args := m.Called(ctx, paymentId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*canonical.Payment), args.Error(1)
}
