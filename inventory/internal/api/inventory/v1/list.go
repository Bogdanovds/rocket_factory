package v1

import (
	"context"
	"github.com/bogdanovds/rocket_factory/inventory/internal/converter"

	inventoryV1 "github.com/bogdanovds/rocket_factory/shared/pkg/proto/inventory/v1"
)

func (a *InventoryAPI) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	modelParts, err := a.partService.ListParts(ctx, converter.ToProtoPartsFilter(req.Filter))
	if err != nil {
		return nil, err
	}

	protoParts := make([]*inventoryV1.Part, 0, len(modelParts))
	for _, part := range modelParts {
		protoParts = append(protoParts, converter.ToProtoPart(part))
	}

	return &inventoryV1.ListPartsResponse{Parts: protoParts}, nil
}
