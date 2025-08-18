package v1

import (
	"context"
	"errors"
	"fmt"
	"github.com/bogdanovds/rocket_factory/order/internal/converter"
	"github.com/bogdanovds/rocket_factory/order/internal/model"
	orderV1 "github.com/bogdanovds/rocket_factory/shared/pkg/openapi/order/v1"
	"github.com/google/uuid"
)

func (h *Handler) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	orderID, err := uuid.Parse(params.OrderUUID.String())
	if err != nil {
		return badRequest("invalid order UUID format"), nil
	}

	order, err := h.service.GetOrder(ctx, orderID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return notFound(fmt.Sprintf("Order with UUID %s not found", params.OrderUUID)), nil
		}
		return nil, fmt.Errorf("get order error: %w", err)
	}

	return converter.ConvertOrderToDTO(order), nil
}
