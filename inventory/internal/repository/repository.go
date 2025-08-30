package repository

import (
	"context"

	"github.com/bogdanovds/rocket_factory/inventory/internal/model"
)

type PartRepository interface {
	Get(ctx context.Context, uuid string) (*model.Part, error)
	List(ctx context.Context) ([]*model.Part, error)
}
