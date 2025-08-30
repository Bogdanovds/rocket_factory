package v1

import (
	"context"
	"fmt"

	"github.com/bogdanovds/rocket_factory/order/internal/model"
	orderV1 "github.com/bogdanovds/rocket_factory/shared/pkg/openapi/order/v1"
	"github.com/google/uuid"
)

func (h *Handler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	orderID, err := uuid.Parse(params.OrderUUID.String())
	if err != nil {
		return badRequest("invalid order UUID format"), nil
	}

	order, err := h.service.PayOrder(ctx, orderID, string(req.PaymentMethod))
	if err != nil {
		switch err {
		case model.ErrOrderNotFound:
			return notFound(fmt.Sprintf("Order with UUID %s not found", params.OrderUUID)), nil
		case model.ErrOrderAlreadyPaid, model.ErrOrderCancelled, model.ErrOrderFulfilled:
			return conflict(err.Error()), nil
		case model.ErrPaymentRequired:
			return badRequest(err.Error()), nil
		default:
			return nil, fmt.Errorf("payment processing error: %w", err)
		}
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: order.TransactionID,
	}, nil
}
