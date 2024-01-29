package service

import (
	"context"
	"errors"
	"net/http"
	"tech-challenge-payment/internal/canonical"
	"tech-challenge-payment/internal/integration/order"
	"tech-challenge-payment/internal/mocks"
	"tech-challenge-payment/internal/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	time := time.Now()
	type Given struct {
		payment     canonical.Payment
		paymentRepo func() repository.PaymentRepository
	}
	type Expected struct {
		err assert.ErrorAssertionFunc
	}
	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"given payment with main fields filled, must return created paymend with all fields filled": {
			given: Given{
				payment: canonical.Payment{
					ID:          canonical.NewUUID(),
					OrderID:     "1234",
					PaymentType: 3,
					CreatedAt:   time,
					UpdatedAt:   time,
					Status:      canonical.PAYMENT_CREATED,
				},
				paymentRepo: func() repository.PaymentRepository {
					payment := canonical.Payment{
						ID:          canonical.NewUUID(),
						OrderID:     "1234",
						PaymentType: 3,
						CreatedAt:   time,
						UpdatedAt:   time,
						Status:      canonical.PAYMENT_CREATED,
					}
					repoMock := &mocks.PaymentRepositoryMock{}
					repoMock.On("Create", mock.Anything, mock.MatchedBy(func(payment canonical.Payment) bool {
						return payment.OrderID == "1234"
					})).Return(payment, nil)
					return repoMock
				},
			},
			expected: Expected{
				err: assert.NoError,
			},
		},
		"given error creating, must return error": {
			given: Given{
				payment: canonical.Payment{
					ID:          canonical.NewUUID(),
					OrderID:     "1234",
					PaymentType: 3,
					CreatedAt:   time,
					UpdatedAt:   time,
					Status:      canonical.PAYMENT_CREATED,
				},
				paymentRepo: func() repository.PaymentRepository {
					repoMock := &mocks.PaymentRepositoryMock{}
					repoMock.On("Create", mock.Anything, mock.Anything).Return(canonical.Payment{
						ID:          canonical.NewUUID(),
						OrderID:     "1234",
						PaymentType: 3,
						CreatedAt:   time,
						UpdatedAt:   time,
						Status:      canonical.PAYMENT_CREATED,
					}, errors.New("error creating payment"))
					return repoMock
				},
			},
			expected: Expected{
				err: assert.Error,
			},
		},
	}

	for _, tc := range tests {
		_, err := NewPaymentService(tc.given.paymentRepo(), order.NewOrderService(http.DefaultClient)).Create(context.Background(), tc.given.payment)

		tc.expected.err(t, err)
	}
}
func TestGetByID(t *testing.T) {

	type Given struct {
		id          string
		paymentRepo func() repository.PaymentRepository
	}
	type Expected struct {
		err assert.ErrorAssertionFunc
	}
	tests := map[string]struct {
		given    Given
		expected Expected
	}{

		"given payment with main fields filled, must return created paymend with all fields filled": {
			given: Given{
				id: "1234",
				paymentRepo: func() repository.PaymentRepository {
					repoMock := &mocks.PaymentRepositoryMock{}
					repoMock.On("GetByID", mock.Anything, "1234").Return(&canonical.Payment{
						ID:          canonical.NewUUID(),
						OrderID:     "1234",
						PaymentType: 3,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
						Status:      canonical.PAYMENT_CREATED,
					}, nil)
					return repoMock
				},
			},
			expected: Expected{
				err: assert.NoError,
			},
		},
	}

	for _, tc := range tests {
		_, err := NewPaymentService(tc.given.paymentRepo(), order.NewOrderService(http.DefaultClient)).GetByID(context.Background(), tc.given.id)

		tc.expected.err(t, err)
	}
}

func TestCallback(t *testing.T) {
	payment := &canonical.Payment{
		ID:          canonical.NewUUID(),
		OrderID:     "1234",
		PaymentType: 3,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Status:      canonical.PAYMENT_CREATED,
	}
	type Given struct {
		id          string
		status      canonical.PaymentStatus
		paymentRepo func() repository.PaymentRepository
	}
	type Expected struct {
		err assert.ErrorAssertionFunc
	}
	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"given payment payment found, update status": {
			given: Given{
				id:     payment.OrderID,
				status: canonical.PAYMENT_FAILED,
				paymentRepo: func() repository.PaymentRepository {
					repoMock := &mocks.PaymentRepositoryMock{}
					repoMock.On("GetByID", mock.Anything, payment.OrderID).Return(payment, nil)
					repoMock.On("Update", mock.Anything, payment.OrderID, mock.MatchedBy(func(input canonical.Payment) bool {
						return input.OrderID == payment.OrderID
					})).Return(nil)
					return repoMock
				},
			},
			expected: Expected{
				err: assert.NoError,
			},
		},
		"given error on db search": {
			given: Given{
				id:     payment.OrderID,
				status: canonical.PAYMENT_FAILED,
				paymentRepo: func() repository.PaymentRepository {
					repoMock := &mocks.PaymentRepositoryMock{}
					repoMock.On("GetByID", mock.Anything, payment.OrderID).Return(nil, errors.New("db error"))
					return repoMock
				},
			},
			expected: Expected{
				err: assert.Error,
			},
		},
		"given payment not found": {
			given: Given{
				id:     payment.OrderID,
				status: canonical.PAYMENT_FAILED,
				paymentRepo: func() repository.PaymentRepository {
					repoMock := &mocks.PaymentRepositoryMock{}
					repoMock.On("GetByID", mock.Anything, payment.OrderID).Return(nil, nil)
					return repoMock
				},
			},
			expected: Expected{
				err: assert.Error,
			},
		},
		"given update error return error": {
			given: Given{
				id:     payment.OrderID,
				status: canonical.PAYMENT_FAILED,
				paymentRepo: func() repository.PaymentRepository {
					repoMock := &mocks.PaymentRepositoryMock{}

					repoMock.On("GetByID", mock.Anything, payment.OrderID).Return(payment, nil)
					repoMock.On("Update", mock.Anything, payment.OrderID, mock.MatchedBy(func(input canonical.Payment) bool {
						return input.OrderID == payment.OrderID
					})).Return(errors.New("db error"))
					return repoMock
				},
			},
			expected: Expected{
				err: assert.Error,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			orderMock := &order.OrderServiceMock{}
			orderMock.On("UpdateStatus", payment.OrderID, "CANCELLED").Return(nil)
			err := NewPaymentService(tc.given.paymentRepo(), orderMock).Callback(context.Background(), tc.given.id, tc.given.status)

			tc.expected.err(t, err)
		})
	}
}

func setupRepoMock(payment canonical.Payment) *mocks.PaymentRepositoryMock {
	repoMock := &mocks.PaymentRepositoryMock{}
	repoMock.On("Create", mock.Anything, payment).Return(payment, nil)
	return repoMock
}
