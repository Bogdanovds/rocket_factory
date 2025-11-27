package app

import (
	"context"

	api "github.com/bogdanovds/rocket_factory/payment/internal/api/payment/v1"
	"github.com/bogdanovds/rocket_factory/payment/internal/service"
	paymentService "github.com/bogdanovds/rocket_factory/payment/internal/service/payment"
	paymentV1 "github.com/bogdanovds/rocket_factory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentV1API paymentV1.PaymentServiceServer

	paymentService service.PaymentService
}

func newDIContainer() *diContainer {
	return &diContainer{}
}

// PaymentV1API возвращает gRPC API сервер
func (d *diContainer) PaymentV1API(ctx context.Context) paymentV1.PaymentServiceServer {
	if d.paymentV1API == nil {
		d.paymentV1API = api.NewPaymentAPI(d.PaymentService(ctx))
	}

	return d.paymentV1API
}

// PaymentService возвращает сервис оплаты
func (d *diContainer) PaymentService(ctx context.Context) service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = paymentService.NewPaymentService()
	}

	return d.paymentService
}

