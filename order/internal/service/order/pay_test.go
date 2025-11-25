package order

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/bogdanovds/rocket_factory/order/internal/model"
)

func (s *OrderServiceTestSuite) TestPayOrder_Success() {
	ctx := context.Background()
	orderID := uuid.New()
	userID := uuid.New()
	transactionID := uuid.New()
	paymentMethod := "CARD"

	existingOrder := &model.Order{
		ID:         orderID,
		UserID:     userID,
		PartIDs:    []uuid.UUID{uuid.New()},
		TotalPrice: 150.0,
		Status:     model.OrderStatusPending,
	}

	s.mockRepo.On("Get", ctx, orderID).Return(existingOrder, nil)
	s.mockPaymentClient.On("PayOrder", ctx, orderID, userID, paymentMethod).Return(transactionID, nil)
	s.mockRepo.On("Update", ctx, mock.AnythingOfType("*model.Order")).Return(nil)

	order, err := s.service.PayOrder(ctx, orderID, paymentMethod)

	s.NoError(err)
	s.NotNil(order)
	s.Equal(model.OrderStatusPaid, order.Status)
	s.Equal(transactionID, order.TransactionID)
	s.Equal(paymentMethod, order.PaymentMethod)
}

func (s *OrderServiceTestSuite) TestPayOrder_OrderNotFound() {
	ctx := context.Background()
	orderID := uuid.New()

	s.mockRepo.On("Get", ctx, orderID).Return(nil, model.ErrOrderNotFound)

	order, err := s.service.PayOrder(ctx, orderID, "CARD")

	s.Nil(order)
	s.ErrorIs(err, model.ErrOrderNotFound)
}

func (s *OrderServiceTestSuite) TestPayOrder_AlreadyPaid() {
	ctx := context.Background()
	orderID := uuid.New()

	existingOrder := &model.Order{
		ID:     orderID,
		Status: model.OrderStatusPaid,
	}

	s.mockRepo.On("Get", ctx, orderID).Return(existingOrder, nil)

	order, err := s.service.PayOrder(ctx, orderID, "CARD")

	s.Nil(order)
	s.ErrorIs(err, model.ErrOrderAlreadyPaid)
}

func (s *OrderServiceTestSuite) TestPayOrder_OrderCancelled() {
	ctx := context.Background()
	orderID := uuid.New()

	existingOrder := &model.Order{
		ID:     orderID,
		Status: model.OrderStatusCancelled,
	}

	s.mockRepo.On("Get", ctx, orderID).Return(existingOrder, nil)

	order, err := s.service.PayOrder(ctx, orderID, "CARD")

	s.Nil(order)
	s.ErrorIs(err, model.ErrOrderCancelled)
}

func (s *OrderServiceTestSuite) TestPayOrder_OrderFulfilled() {
	ctx := context.Background()
	orderID := uuid.New()

	existingOrder := &model.Order{
		ID:     orderID,
		Status: model.OrderStatusFulfilled,
	}

	s.mockRepo.On("Get", ctx, orderID).Return(existingOrder, nil)

	order, err := s.service.PayOrder(ctx, orderID, "CARD")

	s.Nil(order)
	s.ErrorIs(err, model.ErrOrderFulfilled)
}

func (s *OrderServiceTestSuite) TestPayOrder_EmptyPaymentMethod() {
	ctx := context.Background()
	orderID := uuid.New()

	existingOrder := &model.Order{
		ID:     orderID,
		Status: model.OrderStatusPending,
	}

	s.mockRepo.On("Get", ctx, orderID).Return(existingOrder, nil)

	order, err := s.service.PayOrder(ctx, orderID, "")

	s.Nil(order)
	s.ErrorIs(err, model.ErrPaymentRequired)
}

func (s *OrderServiceTestSuite) TestPayOrder_PaymentFailed() {
	ctx := context.Background()
	orderID := uuid.New()
	userID := uuid.New()

	existingOrder := &model.Order{
		ID:     orderID,
		UserID: userID,
		Status: model.OrderStatusPending,
	}

	s.mockRepo.On("Get", ctx, orderID).Return(existingOrder, nil)
	s.mockPaymentClient.On("PayOrder", ctx, orderID, userID, "CARD").Return(uuid.Nil, errors.New("payment failed"))

	order, err := s.service.PayOrder(ctx, orderID, "CARD")

	s.Nil(order)
	s.Error(err)
	s.Contains(err.Error(), "payment failed")
}

