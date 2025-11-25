package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/bogdanovds/rocket_factory/order/internal/model"
)

type Service interface {
	CreateOrder(ctx context.Context, userID uuid.UUID, partIDs []uuid.UUID) (*model.Order, error)
	GetOrder(ctx context.Context, orderID uuid.UUID) (*model.Order, error)
	PayOrder(ctx context.Context, orderID uuid.UUID, paymentMethod string) (*model.Order, error)
	CancelOrder(ctx context.Context, orderID uuid.UUID) error
}
