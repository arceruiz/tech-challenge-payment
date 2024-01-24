package rest

import "tech-challenge-payment/internal/canonical"

func (pr *PaymentRequest) toCanonical() canonical.Payment {
	return canonical.Payment{
		PaymentType: pr.PaymentType,
		CreatedAt:   pr.CreatedAt,
		UpdatedAt:   pr.UpdatedAt,
		Status:      canonical.PaymentStatus(pr.Status),
		OrderID:     pr.OrderID,
	}
}
