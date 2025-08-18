package converter

import (
	"github.com/bogdanovds/rocket_factory/order/internal/model"
	inventoryV1 "github.com/bogdanovds/rocket_factory/shared/pkg/proto/inventory/v1"
	"github.com/google/uuid"
)

func convertProtoToPart(part *inventoryV1.Part) *model.Part {
	return &model.Part{
		ID:       mustParseUUID(part.Uuid),
		Name:     part.Name,
		Price:    float64(part.Price),
		Category: string(part.Category),
	}
}

func mustParseUUID(str string) uuid.UUID {
	id, err := uuid.Parse(str)
	if err != nil {
		panic(err)
	}
	return id
}
