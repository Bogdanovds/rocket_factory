package inventory

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"github.com/bogdanovds/rocket_factory/order/internal/model"
	inventoryV1 "github.com/bogdanovds/rocket_factory/shared/pkg/proto/inventory/v1"
)

type Client struct {
	client inventoryV1.InventoryServiceClient
}

func New(conn *grpc.ClientConn) *Client {
	return &Client{
		client: inventoryV1.NewInventoryServiceClient(conn),
	}
}

func (c *Client) ListParts(ctx context.Context, partIDs []uuid.UUID) ([]*model.Part, error) {
	partUUIDs := make([]string, len(partIDs))
	for i, id := range partIDs {
		partUUIDs[i] = id.String()
	}

	resp, err := c.client.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Uuids: partUUIDs,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("gRPC inventory error: %w", err)
	}

	parts := make([]*model.Part, len(resp.Parts))
	for i, p := range resp.Parts {
		parts[i] = convertProtoToPart(p)
	}

	return parts, nil
}

func convertProtoToPart(part *inventoryV1.Part) *model.Part {
	id, err := uuid.Parse(part.Uuid)
	if err != nil {
		// В реальном проекте нужно обработать ошибку
		id = uuid.Nil
	}

	return &model.Part{
		ID:       id,
		Name:     part.Name,
		Price:    float64(part.Price),
		Category: string(part.Category),
	}
}
