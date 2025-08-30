package order

import (
	"sync"

	"github.com/bogdanovds/rocket_factory/order/internal/model"
	"github.com/google/uuid"
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
