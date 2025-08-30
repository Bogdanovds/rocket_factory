package v1

import (
	"github.com/bogdanovds/rocket_factory/inventory/internal/service"
	inventoryV1 "github.com/bogdanovds/rocket_factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
)

type InventoryAPI struct {
	inventoryV1.UnimplementedInventoryServiceServer
	partService service.PartService
}

func NewInventoryAPI(partService service.PartService) *InventoryAPI {
	return &InventoryAPI{
		partService: partService,
	}
}

func RegisterInventoryServiceServer(s *grpc.Server, api *InventoryAPI) {
	inventoryV1.RegisterInventoryServiceServer(s, api)
}
