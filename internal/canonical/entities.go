package canonical

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrorNotFound = fmt.Errorf("entity not found")
)

type Payment struct {
	ID          string        `bson:"_id"`
	OrderID     string        `bson:"order_id"`
	PaymentType int           `bson:"payment_type"`
	CreatedAt   time.Time     `bson:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at"`
	Status      PaymentStatus `bson:"status"`
}

type PaymentStatus int

const (
	PAYMENT_CREATED PaymentStatus = iota
	PAYMENT_PAYED
	PAYMENT_FAILED
)

var MapPaymentStatus = map[string]PaymentStatus{
	"OK":        PAYMENT_PAYED,
	"NOK":       PAYMENT_FAILED,
	"ERROR":     PAYMENT_FAILED,
	"INIT":      PAYMENT_CREATED,
	"":          PAYMENT_FAILED,
	"COMPLETED": PAYMENT_PAYED,
	"PENDING":   PAYMENT_CREATED,
}

func NewUUID() string {
	return uuid.New().String()
}

func HandleError(err error) error {
	if errors.Is(err, ErrorNotFound) {
		return err
	}
	return fmt.Errorf("unexpected error occurred %w", err)

}
