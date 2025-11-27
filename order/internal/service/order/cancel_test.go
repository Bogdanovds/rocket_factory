package order

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/bogdanovds/rocket_factory/order/internal/model"
)

func (s *OrderServiceTestSuite) TestCancelOrder_Success() {
	ctx := context.Background()
	orderID := uuid.New()

	existingOrder := &model.Order{
		ID:     orderID,
		Status: model.OrderStatusPending,
	}

	s.mockRepo.On("Get", ctx, orderID).Return(existingOrder, nil)
	s.mockRepo.On("Update", ctx, mock.AnythingOfType("*model.Order")).Return(nil)

	err := s.service.CancelOrder(ctx, orderID)

	s.NoError(err)
}

func (s *OrderServiceTestSuite) TestCancelOrder_OrderNotFound() {
	ctx := context.Background()
	orderID := uuid.New()

	s.mockRepo.On("Get", ctx, orderID).Return(nil, model.ErrOrderNotFound)

	err := s.service.CancelOrder(ctx, orderID)

	s.ErrorIs(err, model.ErrOrderNotFound)
}

func (s *OrderServiceTestSuite) TestCancelOrder_AlreadyPaid() {
	ctx := context.Background()
	orderID := uuid.New()

	existingOrder := &model.Order{
		ID:     orderID,
		Status: model.OrderStatusPaid,
	}

	s.mockRepo.On("Get", ctx, orderID).Return(existingOrder, nil)

	err := s.service.CancelOrder(ctx, orderID)

	s.ErrorIs(err, model.ErrOrderAlreadyPaid)
}

func (s *OrderServiceTestSuite) TestCancelOrder_AlreadyCancelled() {
	ctx := context.Background()
	orderID := uuid.New()

	existingOrder := &model.Order{
		ID:     orderID,
		Status: model.OrderStatusCancelled,
	}

	s.mockRepo.On("Get", ctx, orderID).Return(existingOrder, nil)

	err := s.service.CancelOrder(ctx, orderID)

	s.ErrorIs(err, model.ErrOrderCancelled)
}

func (s *OrderServiceTestSuite) TestCancelOrder_OrderFulfilled() {
	ctx := context.Background()
	orderID := uuid.New()

	existingOrder := &model.Order{
		ID:     orderID,
		Status: model.OrderStatusFulfilled,
	}

	s.mockRepo.On("Get", ctx, orderID).Return(existingOrder, nil)

	err := s.service.CancelOrder(ctx, orderID)

	s.ErrorIs(err, model.ErrOrderFulfilled)
}

func (s *OrderServiceTestSuite) TestCancelOrder_RepositoryError() {
	ctx := context.Background()
	orderID := uuid.New()

	existingOrder := &model.Order{
		ID:     orderID,
		Status: model.OrderStatusPending,
	}

	s.mockRepo.On("Get", ctx, orderID).Return(existingOrder, nil)
	s.mockRepo.On("Update", ctx, mock.AnythingOfType("*model.Order")).Return(errors.New("db error"))

	err := s.service.CancelOrder(ctx, orderID)

	s.Error(err)
	s.Contains(err.Error(), "repository error")
}
