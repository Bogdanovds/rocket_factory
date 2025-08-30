package order

import (
	"context"

	"github.com/bogdanovds/rocket_factory/order/internal/model"
)

func (r *Repository) Update(ctx context.Context, order *model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[order.ID] = order
	return nil
}
