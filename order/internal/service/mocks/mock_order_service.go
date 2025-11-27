package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/bogdanovds/rocket_factory/order/internal/model"
)

// MockOrderService - мок сервиса заказов
type MockOrderService struct {
	mock.Mock
}

// NewMockOrderService создает новый мок сервиса
func NewMockOrderService() *MockOrderService {
	return &MockOrderService{}
}

// CreateOrder создает новый заказ
func (m *MockOrderService) CreateOrder(ctx context.Context, userID uuid.UUID, partIDs []uuid.UUID) (*model.Order, error) {
	args := m.Called(ctx, userID, partIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

// GetOrder возвращает заказ по ID
func (m *MockOrderService) GetOrder(ctx context.Context, orderID uuid.UUID) (*model.Order, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

// PayOrder оплачивает заказ
func (m *MockOrderService) PayOrder(ctx context.Context, orderID uuid.UUID, paymentMethod string) (*model.Order, error) {
	args := m.Called(ctx, orderID, paymentMethod)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

// CancelOrder отменяет заказ
func (m *MockOrderService) CancelOrder(ctx context.Context, orderID uuid.UUID) error {
	args := m.Called(ctx, orderID)
	return args.Error(0)
}
