package part

import (
	"context"
	"errors"

	"github.com/bogdanovds/rocket_factory/inventory/internal/model"
)

func (s *PartServiceTestSuite) TestGetPart_Success() {
	ctx := context.Background()
	expectedPart := &model.Part{
		Uuid: "test-uuid-123",
		Name: "Test Part",
	}

	s.mockRepo.On("Get", ctx, "test-uuid-123").Return(expectedPart, nil)

	part, err := s.service.GetPart(ctx, "test-uuid-123")

	s.NoError(err)
	s.NotNil(part)
	s.Equal(expectedPart.Uuid, part.Uuid)
	s.Equal(expectedPart.Name, part.Name)
}

func (s *PartServiceTestSuite) TestGetPart_NotFound() {
	ctx := context.Background()

	s.mockRepo.On("Get", ctx, "non-existent-uuid").Return(nil, model.ErrPartNotFound)

	part, err := s.service.GetPart(ctx, "non-existent-uuid")

	s.Nil(part)
	s.ErrorIs(err, model.ErrPartNotFound)
}

func (s *PartServiceTestSuite) TestGetPart_RepositoryError() {
	ctx := context.Background()
	repoErr := errors.New("database connection failed")

	s.mockRepo.On("Get", ctx, "test-uuid").Return(nil, repoErr)

	part, err := s.service.GetPart(ctx, "test-uuid")

	s.Nil(part)
	s.ErrorIs(err, model.ErrRepositoryOperation)
}
