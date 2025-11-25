package part

import (
	"context"

	"github.com/bogdanovds/rocket_factory/inventory/internal/model"
)

func (r *Repository) List(_ context.Context) ([]*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	parts := make([]*model.Part, 0, len(r.parts))
	for _, part := range r.parts {
		parts = append(parts, part)
	}

	return parts, nil
}
