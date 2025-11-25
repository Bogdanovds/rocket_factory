package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/bogdanovds/rocket_factory/order/internal/model"
)

func (s *Service) CancelOrder(ctx context.Context, orderID uuid.UUID) error {
	order, err := s.repo.Get(ctx, orderID)
	if err != nil {
		return err
	}

	switch order.Status {
	case model.OrderStatusPaid:
		return model.ErrOrderAlreadyPaid
	case model.OrderStatusCancelled:
		return model.ErrOrderCancelled
	case model.OrderStatusFulfilled:
		return model.ErrOrderFulfilled
	}

	order.Status = model.OrderStatusCancelled
	if err := s.repo.Update(ctx, order); err != nil {
		return fmt.Errorf("repository error: %w", err)
	}

	return nil
}
