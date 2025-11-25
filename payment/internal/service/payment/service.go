package payment

import (
	"context"

	paymentV1 "github.com/bogdanovds/rocket_factory/shared/pkg/proto/payment/v1"
)

// Service реализует PaymentService
type Service struct {
	paymentV1.UnimplementedPaymentServiceServer
}

// NewPaymentService создает новый экземпляр сервиса оплаты
func NewPaymentService() *Service {
	return &Service{}
}

// PayOrder обрабатывает оплату заказа
func (s *Service) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	return payOrder(ctx, req)
}
