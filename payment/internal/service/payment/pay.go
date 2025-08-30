package payment

import (
	"context"
	"log"

	"github.com/google/uuid"

	paymentV1 "github.com/bogdanovds/rocket_factory/shared/pkg/proto/payment/v1"
)

// payOrder обрабатывает платеж и возвращает UUID транзакции
func payOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	transactionUUID := uuid.New().String()

	log.Printf("Оплата прошла успешно, transaction_uuid: %s\n"+
		"Детали платежа:\n"+
		" - Order UUID: %s\n"+
		" - User UUID: %s\n"+
		" - Payment Method: %s",
		transactionUUID, req.OrderUuid, req.UserUuid, req.PaymentMethod.String())

	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}
