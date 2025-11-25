package part

import (
	"context"

	"github.com/bogdanovds/rocket_factory/inventory/internal/model"
)

func (r *Repository) Get(_ context.Context, uuid string) (*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, exists := r.parts[uuid]
	if !exists {
		return nil, model.ErrPartNotFound
	}

	return part, nil
}
