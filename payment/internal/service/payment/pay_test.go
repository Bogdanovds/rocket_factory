package payment

import (
	"context"

	"github.com/google/uuid"

	paymentV1 "github.com/bogdanovds/rocket_factory/shared/pkg/proto/payment/v1"
)

func (s *PaymentServiceTestSuite) TestPayOrder_Success_CardPayment() {
	ctx := context.Background()
	orderUUID := uuid.New().String()
	userUUID := uuid.New().String()

	req := &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	resp, err := s.service.PayOrder(ctx, req)

	s.NoError(err)
	s.NotNil(resp)
	s.NotEmpty(resp.TransactionUuid)

	// Проверяем, что transaction_uuid - валидный UUID
	_, parseErr := uuid.Parse(resp.TransactionUuid)
	s.NoError(parseErr)
}

func (s *PaymentServiceTestSuite) TestPayOrder_Success_SBPPayment() {
	ctx := context.Background()
	orderUUID := uuid.New().String()
	userUUID := uuid.New().String()

	req := &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_SBP,
	}

	resp, err := s.service.PayOrder(ctx, req)

	s.NoError(err)
	s.NotNil(resp)
	s.NotEmpty(resp.TransactionUuid)
}

func (s *PaymentServiceTestSuite) TestPayOrder_Success_UnspecifiedPayment() {
	ctx := context.Background()
	orderUUID := uuid.New().String()
	userUUID := uuid.New().String()

	req := &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED,
	}

	resp, err := s.service.PayOrder(ctx, req)

	s.NoError(err)
	s.NotNil(resp)
	s.NotEmpty(resp.TransactionUuid)
}

func (s *PaymentServiceTestSuite) TestPayOrder_UniqueTransactionIDs() {
	ctx := context.Background()
	orderUUID := uuid.New().String()
	userUUID := uuid.New().String()

	req := &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	resp1, err1 := s.service.PayOrder(ctx, req)
	resp2, err2 := s.service.PayOrder(ctx, req)

	s.NoError(err1)
	s.NoError(err2)
	s.NotEqual(resp1.TransactionUuid, resp2.TransactionUuid)
}
