package v1

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/bogdanovds/rocket_factory/order/internal/model"
	orderV1 "github.com/bogdanovds/rocket_factory/shared/pkg/openapi/order/v1"
)

func (h *Handler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	userID, err := uuid.Parse(req.UserUUID.String())
	if err != nil {
		return badRequest("invalid user UUID"), nil
	}

	partIDs := make([]uuid.UUID, len(req.PartUuids))
	for i, partUUID := range req.PartUuids {
		partIDs[i] = uuid.MustParse(partUUID.String())
	}

	order, err := h.service.CreateOrder(ctx, userID, partIDs)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrPartsNotSpecified):
			return badRequest(err.Error()), nil
		case errors.Is(err, model.ErrPartsNotFound):
			return notFound(err.Error()), nil
		default:
			return nil, fmt.Errorf("service error: %w", err)
		}
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  order.ID,
		TotalPrice: float32(order.TotalPrice),
	}, nil
}
