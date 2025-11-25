package service

import (
	"context"

	"github.com/bogdanovds/rocket_factory/inventory/internal/model"
)

type PartService interface {
	GetPart(ctx context.Context, uuid string) (*model.Part, error)
	ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error)
}
