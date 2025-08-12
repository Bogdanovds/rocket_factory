package part

import (
	"context"
	"errors"
	"github.com/bogdanovds/rocket_factory/inventory/internal/model"
)

func (s *Service) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	part, err := s.repo.Get(ctx, uuid)
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, model.ErrPartNotFound
		}
		return nil, model.ErrRepositoryOperation
	}

	return part, nil
}
