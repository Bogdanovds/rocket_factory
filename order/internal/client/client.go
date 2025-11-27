package client

import (
	"context"

	"github.com/google/uuid"

	"github.com/bogdanovds/rocket_factory/order/internal/model"
)

// InventoryClient - интерфейс клиента inventory
type InventoryClient interface {
	ListParts(ctx context.Context, partIDs []uuid.UUID) ([]*model.Part, error)
}

// PaymentClient - интерфейс клиента payment
type PaymentClient interface {
	PayOrder(ctx context.Context, orderID, userID uuid.UUID, method string) (uuid.UUID, error)
}
