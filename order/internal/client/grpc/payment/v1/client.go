package payment

import (
	"context"
	"fmt"

	paymentV1 "github.com/bogdanovds/rocket_factory/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type Client struct {
	client paymentV1.PaymentServiceClient
}

func New(conn *grpc.ClientConn) *Client {
	return &Client{
		client: paymentV1.NewPaymentServiceClient(conn),
	}
}

func (c *Client) PayOrder(ctx context.Context, orderID, userID uuid.UUID, method string) (uuid.UUID, error) {
	protoMethod := convertPaymentMethodToProto(method)
	resp, err := c.client.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     orderID.String(),
		UserUuid:      userID.String(),
		PaymentMethod: protoMethod,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("gRPC payment error: %w", err)
	}

	return uuid.Parse(resp.TransactionUuid)
}

func convertPaymentMethodToProto(method string) paymentV1.PaymentMethod {
	switch method {
	case "CARD":
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case "SBP":
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}
