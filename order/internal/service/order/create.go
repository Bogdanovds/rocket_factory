package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/bogdanovds/rocket_factory/order/internal/model"
)

func (s *Service) CreateOrder(ctx context.Context, userID uuid.UUID, partIDs []uuid.UUID) (*model.Order, error) {
	if len(partIDs) == 0 {
		return nil, model.ErrPartsNotSpecified
	}

	parts, err := s.inventoryClient.ListParts(ctx, partIDs)
	if err != nil {
		return nil, fmt.Errorf("inventory client error: %w", err)
	}

	if len(parts) != len(partIDs) {
		return nil, model.ErrPartsNotFound
	}

	totalPrice := 0.0
	for _, part := range parts {
		totalPrice += part.Price
	}

	order := &model.Order{
		ID:         uuid.New(),
		UserID:     userID,
		PartIDs:    partIDs,
		TotalPrice: totalPrice,
		Status:     model.OrderStatusPending,
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}

	return order, nil
}
