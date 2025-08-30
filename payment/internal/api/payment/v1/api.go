package v1

import (
	"context"

	"github.com/bogdanovds/rocket_factory/payment/internal/service"
	paymentV1 "github.com/bogdanovds/rocket_factory/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
)

// API представляет gRPC API для платежного сервиса
type API struct {
	paymentV1.UnimplementedPaymentServiceServer
	paymentService service.PaymentService
}

// NewPaymentAPI создает новый экземпляр API
func NewPaymentAPI(paymentService service.PaymentService) *API {
	return &API{
		paymentService: paymentService,
	}
}

// RegisterService регистрирует сервис в gRPC сервере
func (a *API) RegisterService(s *grpc.Server) {
	paymentV1.RegisterPaymentServiceServer(s, a)
}

// PayOrder обрабатывает запрос на оплату заказа
func (a *API) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	return a.paymentService.PayOrder(ctx, req)
}
