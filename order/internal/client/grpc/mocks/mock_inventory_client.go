package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/bogdanovds/rocket_factory/order/internal/model"
)

// MockInventoryClient - мок клиента inventory
type MockInventoryClient struct {
	mock.Mock
}

// NewMockInventoryClient создает новый мок клиента
func NewMockInventoryClient() *MockInventoryClient {
	return &MockInventoryClient{}
}

// ListParts возвращает список деталей по ID
func (m *MockInventoryClient) ListParts(ctx context.Context, partIDs []uuid.UUID) ([]*model.Part, error) {
	args := m.Called(ctx, partIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Part), args.Error(1)
}
