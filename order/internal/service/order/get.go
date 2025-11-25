package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/bogdanovds/rocket_factory/order/internal/model"
)

func (s *Service) GetOrder(ctx context.Context, orderID uuid.UUID) (*model.Order, error) {
	return s.repo.Get(ctx, orderID)
}
