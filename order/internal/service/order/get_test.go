package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/bogdanovds/rocket_factory/order/internal/model"
)

func (s *OrderServiceTestSuite) TestGetOrder_Success() {
	ctx := context.Background()
	orderID := uuid.New()
	expectedOrder := &model.Order{
		ID:         orderID,
		UserID:     uuid.New(),
		PartIDs:    []uuid.UUID{uuid.New()},
		TotalPrice: 150.0,
		Status:     model.OrderStatusPending,
	}

	s.mockRepo.On("Get", ctx, orderID).Return(expectedOrder, nil)

	order, err := s.service.GetOrder(ctx, orderID)

	s.NoError(err)
	s.NotNil(order)
	s.Equal(expectedOrder.ID, order.ID)
	s.Equal(expectedOrder.TotalPrice, order.TotalPrice)
}

func (s *OrderServiceTestSuite) TestGetOrder_NotFound() {
	ctx := context.Background()
	orderID := uuid.New()

	s.mockRepo.On("Get", ctx, orderID).Return(nil, model.ErrOrderNotFound)

	order, err := s.service.GetOrder(ctx, orderID)

	s.Nil(order)
	s.ErrorIs(err, model.ErrOrderNotFound)
}
