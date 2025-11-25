package part

import (
	"context"
	"errors"

	"github.com/bogdanovds/rocket_factory/inventory/internal/model"
)

func (s *PartServiceTestSuite) TestListParts_Success_NoFilter() {
	ctx := context.Background()
	expectedParts := []*model.Part{
		{Uuid: "uuid-1", Name: "Part 1", Manufacturer: &model.Manufacturer{Country: "US"}},
		{Uuid: "uuid-2", Name: "Part 2", Manufacturer: &model.Manufacturer{Country: "DE"}},
	}

	s.mockRepo.On("List", ctx).Return(expectedParts, nil)

	parts, err := s.service.ListParts(ctx, nil)

	s.NoError(err)
	s.Len(parts, 2)
}

func (s *PartServiceTestSuite) TestListParts_Success_EmptyFilter() {
	ctx := context.Background()
	expectedParts := []*model.Part{
		{Uuid: "uuid-1", Name: "Part 1", Manufacturer: &model.Manufacturer{Country: "US"}},
	}

	s.mockRepo.On("List", ctx).Return(expectedParts, nil)

	parts, err := s.service.ListParts(ctx, &model.PartsFilter{})

	s.NoError(err)
	s.Len(parts, 1)
}

func (s *PartServiceTestSuite) TestListParts_Success_FilterByUUIDs() {
	ctx := context.Background()
	allParts := []*model.Part{
		{Uuid: "uuid-1", Name: "Part 1", Manufacturer: &model.Manufacturer{Country: "US"}},
		{Uuid: "uuid-2", Name: "Part 2", Manufacturer: &model.Manufacturer{Country: "DE"}},
		{Uuid: "uuid-3", Name: "Part 3", Manufacturer: &model.Manufacturer{Country: "JP"}},
	}

	s.mockRepo.On("List", ctx).Return(allParts, nil)

	filter := &model.PartsFilter{Uuids: []string{"uuid-1", "uuid-3"}}
	parts, err := s.service.ListParts(ctx, filter)

	s.NoError(err)
	s.Len(parts, 2)
	s.Equal("uuid-1", parts[0].Uuid)
	s.Equal("uuid-3", parts[1].Uuid)
}

func (s *PartServiceTestSuite) TestListParts_Success_FilterByCategory() {
	ctx := context.Background()
	allParts := []*model.Part{
		{Uuid: "uuid-1", Name: "Engine Part", Category: model.CategoryEngine, Manufacturer: &model.Manufacturer{Country: "US"}},
		{Uuid: "uuid-2", Name: "Fuel Part", Category: model.CategoryFuel, Manufacturer: &model.Manufacturer{Country: "DE"}},
		{Uuid: "uuid-3", Name: "Wing Part", Category: model.CategoryWing, Manufacturer: &model.Manufacturer{Country: "JP"}},
	}

	s.mockRepo.On("List", ctx).Return(allParts, nil)

	filter := &model.PartsFilter{Categories: []model.Category{model.CategoryEngine}}
	parts, err := s.service.ListParts(ctx, filter)

	s.NoError(err)
	s.Len(parts, 1)
	s.Equal("Engine Part", parts[0].Name)
}

func (s *PartServiceTestSuite) TestListParts_Success_FilterByCountry() {
	ctx := context.Background()
	allParts := []*model.Part{
		{Uuid: "uuid-1", Name: "Part 1", Manufacturer: &model.Manufacturer{Country: "US"}},
		{Uuid: "uuid-2", Name: "Part 2", Manufacturer: &model.Manufacturer{Country: "DE"}},
		{Uuid: "uuid-3", Name: "Part 3", Manufacturer: &model.Manufacturer{Country: "US"}},
	}

	s.mockRepo.On("List", ctx).Return(allParts, nil)

	filter := &model.PartsFilter{ManufacturerCountries: []string{"US"}}
	parts, err := s.service.ListParts(ctx, filter)

	s.NoError(err)
	s.Len(parts, 2)
}

func (s *PartServiceTestSuite) TestListParts_Success_FilterByTags() {
	ctx := context.Background()
	allParts := []*model.Part{
		{Uuid: "uuid-1", Name: "Part 1", Tags: []string{"heavy", "hot"}, Manufacturer: &model.Manufacturer{Country: "US"}},
		{Uuid: "uuid-2", Name: "Part 2", Tags: []string{"light"}, Manufacturer: &model.Manufacturer{Country: "DE"}},
		{Uuid: "uuid-3", Name: "Part 3", Tags: []string{"heavy"}, Manufacturer: &model.Manufacturer{Country: "JP"}},
	}

	s.mockRepo.On("List", ctx).Return(allParts, nil)

	filter := &model.PartsFilter{Tags: []string{"heavy"}}
	parts, err := s.service.ListParts(ctx, filter)

	s.NoError(err)
	s.Len(parts, 2)
}

func (s *PartServiceTestSuite) TestListParts_RepositoryError() {
	ctx := context.Background()
	repoErr := errors.New("database error")

	s.mockRepo.On("List", ctx).Return(nil, repoErr)

	parts, err := s.service.ListParts(ctx, nil)

	s.Nil(parts)
	s.ErrorIs(err, model.ErrRepositoryOperation)
}

func (s *PartServiceTestSuite) TestListParts_EmptyResult() {
	ctx := context.Background()

	s.mockRepo.On("List", ctx).Return([]*model.Part{}, nil)

	parts, err := s.service.ListParts(ctx, nil)

	s.NoError(err)
	s.Empty(parts)
}

