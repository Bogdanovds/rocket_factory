package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/bogdanovds/rocket_factory/inventory/internal/model"
)

// MockPartService - мок сервиса деталей
type MockPartService struct {
	mock.Mock
}

// NewMockPartService создает новый мок сервиса
func NewMockPartService() *MockPartService {
	return &MockPartService{}
}

// GetPart возвращает деталь по UUID
func (m *MockPartService) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Part), args.Error(1)
}

// ListParts возвращает список деталей с учетом фильтра
func (m *MockPartService) ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Part), args.Error(1)
}

