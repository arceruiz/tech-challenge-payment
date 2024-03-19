package service

import (
	"context"
	"errors"
	"tech-challenge-payment/internal/canonical"
	"tech-challenge-payment/internal/integration/sqs_publisher"
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
					repoMock := &PaymentRepositoryMock{}
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
					repoMock := &PaymentRepositoryMock{}
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
		paymentSvc := paymentService{
			repo: tc.given.paymentRepo(),
		}
		_, err := paymentSvc.Create(context.Background(), tc.given.payment)

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
					repoMock := &PaymentRepositoryMock{}
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
		paymentSvc := paymentService{
			repo: tc.given.paymentRepo(),
		}
		_, err := paymentSvc.GetByID(context.Background(), tc.given.id)

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
		publisher   func() sqs_publisher.Publisher
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
					repoMock := &PaymentRepositoryMock{}
					repoMock.On("GetByID", mock.Anything, payment.OrderID).Return(payment, nil)
					repoMock.On("Update", mock.Anything, payment.OrderID, mock.MatchedBy(func(input canonical.Payment) bool {
						return input.OrderID == payment.OrderID
					})).Return(nil)
					return repoMock
				},
				publisher: func() sqs_publisher.Publisher {
					pubMock := new(PublisherMock)

					pubMock.On("SendMessage").Return(nil)

					return pubMock
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
					repoMock := &PaymentRepositoryMock{}
					repoMock.On("GetByID", mock.Anything, payment.OrderID).Return(nil, errors.New("db error"))
					return repoMock
				},
				publisher: func() sqs_publisher.Publisher {
					return new(PublisherMock)
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
					repoMock := &PaymentRepositoryMock{}
					repoMock.On("GetByID", mock.Anything, payment.OrderID).Return(nil, nil)
					return repoMock
				},
				publisher: func() sqs_publisher.Publisher {
					return new(PublisherMock)
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
					repoMock := &PaymentRepositoryMock{}

					repoMock.On("GetByID", mock.Anything, payment.OrderID).Return(payment, nil)
					repoMock.On("Update", mock.Anything, payment.OrderID, mock.MatchedBy(func(input canonical.Payment) bool {
						return input.OrderID == payment.OrderID
					})).Return(errors.New("db error"))
					return repoMock
				},
				publisher: func() sqs_publisher.Publisher {
					return new(PublisherMock)
				},
			},
			expected: Expected{
				err: assert.Error,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			paymentSvc := paymentService{
				repo:      tc.given.paymentRepo(),
				publisher: tc.given.publisher(),
			}

			err := paymentSvc.Callback(context.Background(), tc.given.id, tc.given.status)

			tc.expected.err(t, err)
		})
	}
}
