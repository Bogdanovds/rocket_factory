package order

import (
	"sync"

	"github.com/google/uuid"

	"github.com/bogdanovds/rocket_factory/order/internal/model"
)

type Repository struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]*model.Order
}

func NewRepo() *Repository {
	return &Repository{
		orders: make(map[uuid.UUID]*model.Order),
	}
}
