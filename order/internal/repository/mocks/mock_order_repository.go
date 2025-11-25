package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/bogdanovds/rocket_factory/order/internal/model"
)

// MockOrderRepository - мок репозитория заказов
type MockOrderRepository struct {
	mock.Mock
}

// NewMockOrderRepository создает новый мок репозитория
func NewMockOrderRepository() *MockOrderRepository {
	return &MockOrderRepository{}
}

// Create создает новый заказ
func (m *MockOrderRepository) Create(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

// Get возвращает заказ по ID
func (m *MockOrderRepository) Get(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

// Update обновляет заказ
func (m *MockOrderRepository) Update(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

