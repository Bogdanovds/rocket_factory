package model

import "github.com/google/uuid"

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusPaid      OrderStatus = "PAID"
	OrderStatusCancelled OrderStatus = "CANCELLED"
	OrderStatusFulfilled OrderStatus = "FULFILLED"
)

type Order struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	PartIDs       []uuid.UUID
	TotalPrice    float64
	Status        OrderStatus
	PaymentMethod string
	TransactionID uuid.UUID
}
