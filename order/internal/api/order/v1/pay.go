package v1

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/bogdanovds/rocket_factory/order/internal/model"
	orderV1 "github.com/bogdanovds/rocket_factory/shared/pkg/openapi/order/v1"
)

func (h *Handler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	orderID, err := uuid.Parse(params.OrderUUID.String())
	if err != nil {
		return badRequest("invalid order UUID format"), nil
	}

	order, err := h.service.PayOrder(ctx, orderID, string(req.PaymentMethod))
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			return notFound(fmt.Sprintf("Order with UUID %s not found", params.OrderUUID)), nil
		case errors.Is(err, model.ErrOrderAlreadyPaid), errors.Is(err, model.ErrOrderCancelled), errors.Is(err, model.ErrOrderFulfilled):
			return conflict(err.Error()), nil
		case errors.Is(err, model.ErrPaymentRequired):
			return badRequest(err.Error()), nil
		default:
			return nil, fmt.Errorf("payment processing error: %w", err)
		}
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: order.TransactionID,
	}, nil
}
