package v1

import (
	"github.com/bogdanovds/rocket_factory/payment/internal/model"
	paymentV1 "github.com/bogdanovds/rocket_factory/shared/pkg/proto/payment/v1"
)

// PayOrderRequestWrapper обертка для запроса оплаты
type PayOrderRequestWrapper struct {
	*paymentV1.PayOrderRequest
}

// Validate валидирует запрос оплаты
func (r *PayOrderRequestWrapper) Validate() error {
	// Здесь можно добавить валидацию полей запроса
	if r.OrderUuid == "" {
		return model.ErrInvalidOrderUUID
	}
	if r.UserUuid == "" {
		return model.ErrInvalidUserUUID
	}
	if r.PaymentMethod == paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED {
		return model.ErrInvalidPaymentMethod
	}
	return nil
}
