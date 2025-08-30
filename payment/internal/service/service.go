package service

import (
	"context"

	paymentV1 "github.com/bogdanovds/rocket_factory/shared/pkg/proto/payment/v1"
)

// PaymentService интерфейс для сервиса оплаты
type PaymentService interface {
	PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error)
}
