package repository

import (
	"context"
	"reflect"
	"tech-challenge-payment/internal/canonical"
	"tech-challenge-payment/internal/config"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/undefinedlabs/go-mpatch"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestGetByID(t *testing.T) {

	mpatch.PatchMethod(time.Now, func() time.Time {
		return time.Date(2020, 11, 01, 00, 00, 00, 0, time.UTC)
	})

	type Given struct {
		mtestFunc func(mt *mtest.T)
	}
	type Expected struct {
	}
	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"given valid search result, must return valid payment": {
			given: Given{
				mtestFunc: func(mt *mtest.T) {
					repo := paymentRepository{
						collection: mt.DB.Collection("fake-collection"),
					}
					mt.AddMockResponses(mtest.CreateCursorResponse(1, "payment.payment", mtest.FirstBatch, bson.D{
						{Key: "_id", Value: "payment_valid"},
						{Key: "order_id", Value: "order_valid"},
						{Key: "payment_type", Value: 0},
						{Key: "created_at", Value: time.Now()},
						{Key: "updated_at", Value: time.Now()},
						{Key: "status", Value: 0},
					}))
					payment, err := repo.GetByID(context.Background(), "payment_valid")
					assert.Nil(t, err)
					assert.Equal(t, payment.ID, "payment_valid")
					assert.Equal(t, payment.OrderID, "order_valid")
					assert.Equal(t, payment.PaymentType, 0)
					assert.Equal(t, payment.CreatedAt, time.Now())
					assert.Equal(t, payment.UpdatedAt, time.Now())
					assert.Equal(t, payment.Status, canonical.PAYMENT_CREATED)
				},
			},
		},
		"given entity not found must return error": {
			given: Given{
				mtestFunc: func(mt *mtest.T) {
					repo := paymentRepository{
						collection: mt.DB.Collection("fake-collection"),
					}
					mt.AddMockResponses(mtest.CreateCursorResponse(0, "payment.payment", mtest.FirstBatch))
					payment, err := repo.GetByID(context.Background(), "asd")
					assert.NotNil(t, err)
					assert.Equal(t, err.Error(), "mongo: no documents in result")
					assert.Nil(t, payment)
				},
			},
		},
	}

	for _, tc := range tests {
		db := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		db.Run("", tc.given.mtestFunc)
	}
}

func TestGetAll(t *testing.T) {

	mpatch.PatchMethod(time.Now, func() time.Time {
		return time.Date(2020, 11, 01, 00, 00, 00, 0, time.UTC)
	})

	type Given struct {
		mtestFunc func(mt *mtest.T)
	}
	type Expected struct {
	}
	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"given valid search result, must return valid payment": {
			given: Given{
				mtestFunc: func(mt *mtest.T) {
					repo := paymentRepository{
						collection: mt.DB.Collection("fake-collection"),
					}
					first := mtest.CreateCursorResponse(1, "payment.payment", mtest.FirstBatch, bson.D{
						{Key: "_id", Value: "payment_valid"},
						{Key: "order_id", Value: "order_valid"},
						{Key: "payment_type", Value: 0},
						{Key: "created_at", Value: time.Now()},
						{Key: "updated_at", Value: time.Now()},
						{Key: "status", Value: 0},
					})
					getMore := mtest.CreateCursorResponse(1, "payment.payment", mtest.NextBatch, bson.D{
						{Key: "_id", Value: "payment_valid"},
						{Key: "order_id", Value: "order_valid"},
						{Key: "payment_type", Value: 0},
						{Key: "created_at", Value: time.Now()},
						{Key: "updated_at", Value: time.Now()},
						{Key: "status", Value: 0},
					})

					lastCursor := mtest.CreateCursorResponse(0, "payment.payment", mtest.NextBatch)
					mt.AddMockResponses(first, getMore, lastCursor)
					payments, err := repo.GetAll(context.Background())
					assert.Nil(t, err)
					for _, payment := range payments {
						assert.Equal(t, payment.ID, "payment_valid")
						assert.Equal(t, payment.OrderID, "order_valid")
						assert.Equal(t, payment.PaymentType, 0)
						assert.Equal(t, payment.CreatedAt, time.Now())
						assert.Equal(t, payment.UpdatedAt, time.Now())
						assert.Equal(t, payment.Status, canonical.PAYMENT_CREATED)
					}
				},
			},
		},
		"given entity not found must return error": {
			given: Given{
				mtestFunc: func(mt *mtest.T) {
					repo := paymentRepository{
						collection: mt.DB.Collection("fake-collection"),
					}
					mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{Message: "mongo: no documents in result"}))
					payment, err := repo.GetAll(context.Background())
					assert.NotNil(t, err)
					assert.Equal(t, err.Error(), "write command error: [{write errors: [{mongo: no documents in result}]}, {<nil>}]")
					assert.Nil(t, payment)
				},
			},
		},
	}

	for _, tc := range tests {
		db := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		db.Run("", tc.given.mtestFunc)
	}
}

func TestCreate(t *testing.T) {

	mpatch.PatchMethod(time.Now, func() time.Time {
		return time.Date(2020, 11, 01, 00, 00, 00, 0, time.UTC)
	})

	type Given struct {
		mtestFunc func(mt *mtest.T)
	}
	type Expected struct {
	}
	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"given given no error saving must return correct entity": {
			given: Given{
				mtestFunc: func(mt *mtest.T) {
					repo := paymentRepository{
						collection: mt.DB.Collection("fake-collection"),
					}
					mt.AddMockResponses(mtest.CreateSuccessResponse())

					tPayment := canonical.Payment{
						ID:          "payment_valid",
						OrderID:     "order_valid",
						PaymentType: 0,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
						Status:      canonical.PAYMENT_CREATED,
					}

					payment, err := repo.Create(context.Background(), tPayment)

					assert.Nil(t, err)
					assert.Equal(t, payment, tPayment)

				},
			},
		},
		"given given error saving must return error": {
			given: Given{
				mtestFunc: func(mt *mtest.T) {
					repo := paymentRepository{
						collection: mt.DB.Collection("fake-collection"),
					}
					mt.AddMockResponses(
						bson.D{
							{Key: "ok", Value: -1},
						},
					)

					tPayment := canonical.Payment{
						ID:          "payment_valid",
						OrderID:     "order_valid",
						PaymentType: 0,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
						Status:      canonical.PAYMENT_CREATED,
					}

					payment, err := repo.Create(context.Background(), tPayment)

					assert.NotNil(t, err)
					assert.Equal(t, payment, tPayment)

				},
			},
		},
	}

	for _, tc := range tests {
		db := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		db.Run("", tc.given.mtestFunc)
	}
}

func TestUpdate(t *testing.T) {

	mpatch.PatchMethod(time.Now, func() time.Time {
		return time.Date(2020, 11, 01, 00, 00, 00, 0, time.UTC)
	})

	type Given struct {
		mtestFunc func(mt *mtest.T)
	}
	type Expected struct {
	}
	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"given given no error updating must return no error": {
			given: Given{
				mtestFunc: func(mt *mtest.T) {
					repo := paymentRepository{
						collection: mt.DB.Collection("fake-collection"),
					}
					mt.AddMockResponses(bson.D{
						{Key: "ok", Value: 1},
						{Key: "value", Value: bson.D{
							{Key: "_id", Value: "payment_valid"},
							{Key: "order_id", Value: "order_valid"},
							{Key: "payment_type", Value: 0},
							{Key: "created_at", Value: time.Now()},
							{Key: "updated_at", Value: time.Now()},
							{Key: "status", Value: 0},
						}},
					})

					tPayment := canonical.Payment{
						ID:          "payment_valid",
						OrderID:     "order_valid",
						PaymentType: 0,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
						Status:      canonical.PAYMENT_CREATED,
					}

					err := repo.Update(context.Background(), "payment_valid", tPayment)

					assert.Nil(t, err)

				},
			},
		},
		"given error saving must return error": {
			given: Given{
				mtestFunc: func(mt *mtest.T) {
					repo := paymentRepository{
						collection: mt.DB.Collection("fake-collection"),
					}
					mt.AddMockResponses(
						bson.D{
							{Key: "ok", Value: -1},
						},
					)
					tPayment := canonical.Payment{
						ID:          "payment_valid",
						OrderID:     "order_valid",
						PaymentType: 0,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
						Status:      canonical.PAYMENT_CREATED,
					}

					err := repo.Update(context.Background(), "payment_valid", tPayment)

					assert.NotNil(t, err)

				},
			},
		},
	}

	for _, tc := range tests {
		db := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		db.Run("", tc.given.mtestFunc)
	}
}

func TestNewMongo(t *testing.T) {

	config.ParseFromFlags()
	got := NewMongo()
	assert.True(t, got != nil && reflect.TypeOf(got).String() == "*mongo.Database")
}
