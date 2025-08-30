package v1

import (
	"context"
	"fmt"

	"github.com/bogdanovds/rocket_factory/order/internal/model"
	orderV1 "github.com/bogdanovds/rocket_factory/shared/pkg/openapi/order/v1"
	"github.com/google/uuid"
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
		switch err {
		case model.ErrPartsNotSpecified:
			return badRequest(err.Error()), nil
		case model.ErrPartsNotFound:
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
