package part

import (
	"context"
	"github.com/bogdanovds/rocket_factory/inventory/internal/model"
	"slices"
)

func (s *Service) ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	parts, err := s.repo.List(ctx)
	if err != nil {
		return nil, model.ErrRepositoryOperation
	}

	if isEmptyFilter(filter) {
		return parts, nil
	}

	var filteredParts []*model.Part
	for _, part := range parts {
		if matchesFilter(part, filter) {
			filteredParts = append(filteredParts, part)
		}
	}

	return filteredParts, nil
}

func isEmptyFilter(filter *model.PartsFilter) bool {
	return filter == nil ||
		(len(filter.Uuids) == 0 &&
			len(filter.Names) == 0 &&
			len(filter.Categories) == 0 &&
			len(filter.ManufacturerCountries) == 0 &&
			len(filter.Tags) == 0)
}

func matchesFilter(part *model.Part, filter *model.PartsFilter) bool {
	if len(filter.Uuids) > 0 && !slices.Contains(filter.Uuids, part.Uuid) {
		return false
	}

	if len(filter.Names) > 0 && !slices.Contains(filter.Names, part.Name) {
		return false
	}

	if len(filter.Categories) > 0 && !containsCategory(filter.Categories, part.Category) {
		return false
	}

	if len(filter.ManufacturerCountries) > 0 && !slices.Contains(filter.ManufacturerCountries, part.Manufacturer.Country) {
		return false
	}

	if len(filter.Tags) > 0 && !hasAnyTag(part.Tags, filter.Tags) {
		return false
	}

	return true
}

func containsCategory(categories []model.Category, category model.Category) bool {
	for _, c := range categories {
		if c == category {
			return true
		}
	}
	return false
}

func hasAnyTag(partTags, filterTags []string) bool {
	for _, ft := range filterTags {
		for _, pt := range partTags {
			if pt == ft {
				return true
			}
		}
	}
	return false
}
