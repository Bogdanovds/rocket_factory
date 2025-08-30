package order

import (
	"context"

	"github.com/bogdanovds/rocket_factory/order/internal/model"
	"github.com/google/uuid"
)

func (s *Service) GetOrder(ctx context.Context, orderID uuid.UUID) (*model.Order, error) {
	return s.repo.Get(ctx, orderID)
}
