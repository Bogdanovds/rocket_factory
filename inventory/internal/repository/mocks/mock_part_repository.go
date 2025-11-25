package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/bogdanovds/rocket_factory/inventory/internal/model"
)

// MockPartRepository - мок репозитория деталей
type MockPartRepository struct {
	mock.Mock
}

// NewMockPartRepository создает новый мок репозитория
func NewMockPartRepository() *MockPartRepository {
	return &MockPartRepository{}
}

// Get возвращает деталь по UUID
func (m *MockPartRepository) Get(ctx context.Context, uuid string) (*model.Part, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Part), args.Error(1)
}

// List возвращает список всех деталей
func (m *MockPartRepository) List(ctx context.Context) ([]*model.Part, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Part), args.Error(1)
}

