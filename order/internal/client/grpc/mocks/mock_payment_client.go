package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockPaymentClient - мок клиента payment
type MockPaymentClient struct {
	mock.Mock
}

// NewMockPaymentClient создает новый мок клиента
func NewMockPaymentClient() *MockPaymentClient {
	return &MockPaymentClient{}
}

// PayOrder оплачивает заказ
func (m *MockPaymentClient) PayOrder(ctx context.Context, orderID, userID uuid.UUID, method string) (uuid.UUID, error) {
	args := m.Called(ctx, orderID, userID, method)
	return args.Get(0).(uuid.UUID), args.Error(1)
}
