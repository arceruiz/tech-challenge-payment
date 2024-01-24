package service_test

import (
	"context"
	"errors"
	"tech-challenge-payment/internal/canonical"
	"tech-challenge-payment/internal/repository"
	repository_test "tech-challenge-payment/internal/repository"
	"tech-challenge-payment/internal/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/undefinedlabs/go-mpatch"
)

func TestCreate(t *testing.T) {

	mpatch.PatchMethod(time.Now, func() time.Time {
		return time.Date(2020, 11, 01, 00, 00, 00, 0, time.UTC)
	})
	mpatch.PatchMethod(canonical.NewUUID, func() string {
		return "4c0fe202-f3e0-4e5f-9868-a2bf1dfb383b"
	})

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
					OrderID:     "1234",
					PaymentType: 3,
				},
				paymentRepo: func() repository.PaymentRepository {
					payment := canonical.Payment{
						ID:          canonical.NewUUID(),
						OrderID:     "1234",
						PaymentType: 3,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
						Status:      canonical.PAYMENT_CREATED,
					}
					repoMock := &repository.PaymentRepositoryMock{}
					repoMock.On("Create", mock.Anything, payment).Return(payment, nil)
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
					OrderID:     "1234",
					PaymentType: 3,
				},
				paymentRepo: func() repository.PaymentRepository {
					repoMock := &repository_test.PaymentRepositoryMock{}
					repoMock.On("Create", mock.Anything, mock.Anything).Return(canonical.Payment{}, errors.New("error creating payment"))
					return repoMock
				},
			},
			expected: Expected{
				err: assert.Error,
			},
		},
	}

	for _, tc := range tests {
		_, err := service.NewPaymentService(tc.given.paymentRepo()).Create(context.Background(), tc.given.payment)

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
					repoMock := &repository.PaymentRepositoryMock{}
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
		_, err := service.NewPaymentService(tc.given.paymentRepo()).GetByID(context.Background(), tc.given.id)

		tc.expected.err(t, err)
	}
}

func TestCallback(t *testing.T) {

	mpatch.PatchMethod(canonical.NewUUID, func() string {
		return "4c0fe202-f3e0-4e5f-9868-a2bf1dfb383b"
	})
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
				id:     "1234",
				status: canonical.PAYMENT_FAILED,
				paymentRepo: func() repository.PaymentRepository {
					repoMock := &repository.PaymentRepositoryMock{}
					repoMock.On("GetByID", mock.Anything, "1234").Return(&canonical.Payment{
						ID:          canonical.NewUUID(),
						OrderID:     "1234",
						PaymentType: 3,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
						Status:      canonical.PAYMENT_CREATED,
					}, nil)
					repoMock.On("Update", mock.Anything, "1234", canonical.Payment{
						ID:          canonical.NewUUID(),
						OrderID:     "1234",
						PaymentType: 3,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
						Status:      canonical.PAYMENT_FAILED,
					}).Return(nil)
					return repoMock
				},
			},
			expected: Expected{
				err: assert.NoError,
			},
		},
		"given error on db search": {
			given: Given{
				id:     "1234",
				status: canonical.PAYMENT_FAILED,
				paymentRepo: func() repository.PaymentRepository {
					repoMock := &repository.PaymentRepositoryMock{}
					repoMock.On("GetByID", mock.Anything, "1234").Return(nil, errors.New("db error"))
					return repoMock
				},
			},
			expected: Expected{
				err: assert.Error,
			},
		},
		"given payment not found": {
			given: Given{
				id:     "1234",
				status: canonical.PAYMENT_FAILED,
				paymentRepo: func() repository.PaymentRepository {
					repoMock := &repository.PaymentRepositoryMock{}
					repoMock.On("GetByID", mock.Anything, "1234").Return(nil, nil)
					return repoMock
				},
			},
			expected: Expected{
				err: assert.Error,
			},
		},
		"given update error return error": {
			given: Given{
				id:     "1234",
				status: canonical.PAYMENT_FAILED,
				paymentRepo: func() repository.PaymentRepository {
					repoMock := &repository.PaymentRepositoryMock{}
					repoMock.On("GetByID", mock.Anything, "1234").Return(&canonical.Payment{
						ID:          canonical.NewUUID(),
						OrderID:     "1234",
						PaymentType: 3,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
						Status:      canonical.PAYMENT_CREATED,
					}, nil)
					repoMock.On("Update", mock.Anything, "1234", canonical.Payment{
						ID:          canonical.NewUUID(),
						OrderID:     "1234",
						PaymentType: 3,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
						Status:      canonical.PAYMENT_FAILED,
					}).Return(errors.New("db error"))
					return repoMock
				},
			},
			expected: Expected{
				err: assert.Error,
			},
		},
	}

	for _, tc := range tests {
		err := service.NewPaymentService(tc.given.paymentRepo()).Callback(context.Background(), tc.given.id, tc.given.status)

		tc.expected.err(t, err)
	}
}

func setupRepoMock(payment canonical.Payment) *repository.PaymentRepositoryMock {
	repoMock := &repository.PaymentRepositoryMock{}
	repoMock.On("Create", mock.Anything, payment).Return(payment, nil)
	return repoMock
}
