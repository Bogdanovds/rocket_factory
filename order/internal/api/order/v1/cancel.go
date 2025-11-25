package v1

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/bogdanovds/rocket_factory/order/internal/model"
	orderV1 "github.com/bogdanovds/rocket_factory/shared/pkg/openapi/order/v1"
)

func (h *Handler) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	orderID, err := uuid.Parse(params.OrderUUID.String())
	if err != nil {
		return badRequest("invalid order UUID format"), nil
	}

	err = h.service.CancelOrder(ctx, orderID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			return notFound(fmt.Sprintf("Order with UUID %s not found", params.OrderUUID)), nil
		case errors.Is(err, model.ErrOrderAlreadyPaid), errors.Is(err, model.ErrOrderCancelled), errors.Is(err, model.ErrOrderFulfilled):
			return conflict(err.Error()), nil
		default:
			return nil, fmt.Errorf("cancel order error: %w", err)
		}
	}

	return &orderV1.CancelOrderNoContent{}, nil
}
