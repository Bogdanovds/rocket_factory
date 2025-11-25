package order

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/bogdanovds/rocket_factory/order/internal/model"
)

func (s *OrderServiceTestSuite) TestCreateOrder_Success() {
	ctx := context.Background()
	userID := uuid.New()
	partIDs := []uuid.UUID{uuid.New(), uuid.New()}

	parts := []*model.Part{
		{ID: partIDs[0], Name: "Part 1", Price: 100.0},
		{ID: partIDs[1], Name: "Part 2", Price: 200.0},
	}

	s.mockInventoryClient.On("ListParts", ctx, partIDs).Return(parts, nil)
	s.mockRepo.On("Create", ctx, mock.AnythingOfType("*model.Order")).Return(nil)

	order, err := s.service.CreateOrder(ctx, userID, partIDs)

	s.NoError(err)
	s.NotNil(order)
	s.Equal(userID, order.UserID)
	s.Equal(partIDs, order.PartIDs)
	s.Equal(300.0, order.TotalPrice)
	s.Equal(model.OrderStatusPending, order.Status)
}

func (s *OrderServiceTestSuite) TestCreateOrder_EmptyParts() {
	ctx := context.Background()
	userID := uuid.New()

	order, err := s.service.CreateOrder(ctx, userID, []uuid.UUID{})

	s.Nil(order)
	s.ErrorIs(err, model.ErrPartsNotSpecified)
}

func (s *OrderServiceTestSuite) TestCreateOrder_InventoryError() {
	ctx := context.Background()
	userID := uuid.New()
	partIDs := []uuid.UUID{uuid.New()}

	s.mockInventoryClient.On("ListParts", ctx, partIDs).Return(nil, errors.New("inventory error"))

	order, err := s.service.CreateOrder(ctx, userID, partIDs)

	s.Nil(order)
	s.Error(err)
	s.Contains(err.Error(), "inventory client error")
}

func (s *OrderServiceTestSuite) TestCreateOrder_PartsNotFound() {
	ctx := context.Background()
	userID := uuid.New()
	partIDs := []uuid.UUID{uuid.New(), uuid.New()}

	// Возвращаем только одну деталь вместо двух
	parts := []*model.Part{
		{ID: partIDs[0], Name: "Part 1", Price: 100.0},
	}

	s.mockInventoryClient.On("ListParts", ctx, partIDs).Return(parts, nil)

	order, err := s.service.CreateOrder(ctx, userID, partIDs)

	s.Nil(order)
	s.ErrorIs(err, model.ErrPartsNotFound)
}

func (s *OrderServiceTestSuite) TestCreateOrder_RepositoryError() {
	ctx := context.Background()
	userID := uuid.New()
	partIDs := []uuid.UUID{uuid.New()}

	parts := []*model.Part{
		{ID: partIDs[0], Name: "Part 1", Price: 100.0},
	}

	s.mockInventoryClient.On("ListParts", ctx, partIDs).Return(parts, nil)
	s.mockRepo.On("Create", ctx, mock.AnythingOfType("*model.Order")).Return(errors.New("db error"))

	order, err := s.service.CreateOrder(ctx, userID, partIDs)

	s.Nil(order)
	s.Error(err)
	s.Contains(err.Error(), "repository error")
}

