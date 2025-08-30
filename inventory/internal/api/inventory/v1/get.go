package v1

import (
	"context"

	"github.com/bogdanovds/rocket_factory/inventory/internal/converter"
	inventoryV1 "github.com/bogdanovds/rocket_factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *InventoryAPI) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part, err := a.partService.GetPart(ctx, req.GetUuid())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.Uuid)
	}

	return &inventoryV1.GetPartResponse{Part: converter.ToProtoPart(part)}, nil
}
