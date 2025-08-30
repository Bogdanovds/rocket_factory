package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/bogdanovds/rocket_factory/order/internal/model"
)

type Repository interface {
	Create(ctx context.Context, order *model.Order) error
	Get(ctx context.Context, id uuid.UUID) (*model.Order, error)
	Update(ctx context.Context, order *model.Order) error
}
