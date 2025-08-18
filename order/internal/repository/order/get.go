package order

import (
	"context"
	"github.com/bogdanovds/rocket_factory/order/internal/model"
	"github.com/google/uuid"
)

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	order, exists := r.orders[id]
	if !exists {
		return nil, model.ErrOrderNotFound
	}
	return order, nil
}
