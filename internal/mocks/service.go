package mocks

import (
	"context"
	"tech-challenge-payment/internal/canonical"

	"github.com/stretchr/testify/mock"
)

type PaymentServiceMock struct {
	mock.Mock
}

func (m *PaymentServiceMock) Create(ctx context.Context, payment canonical.Payment) (*canonical.Payment, error) {
	args := m.Called(ctx, payment)
	return args.Get(0).(*canonical.Payment), args.Error(1)
}

func (m *PaymentServiceMock) Callback(ctx context.Context, paymentId string, status canonical.PaymentStatus) error {
	args := m.Called(ctx, paymentId, status)
	return args.Error(0)
}

func (m *PaymentServiceMock) GetByID(ctx context.Context, paymentId string) (*canonical.Payment, error) {
	args := m.Called(ctx, paymentId)
	return args.Get(0).(*canonical.Payment), args.Error(1)
}
func (m *PaymentServiceMock) GetAll(ctx context.Context) ([]canonical.Payment, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]canonical.Payment), args.Error(1)
}
