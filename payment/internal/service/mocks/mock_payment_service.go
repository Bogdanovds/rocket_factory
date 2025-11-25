package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	paymentV1 "github.com/bogdanovds/rocket_factory/shared/pkg/proto/payment/v1"
)

// MockPaymentService - мок сервиса оплаты
type MockPaymentService struct {
	mock.Mock
}

// NewMockPaymentService создает новый мок сервиса
func NewMockPaymentService() *MockPaymentService {
	return &MockPaymentService{}
}

// PayOrder обрабатывает оплату заказа
func (m *MockPaymentService) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*paymentV1.PayOrderResponse), args.Error(1)
}

