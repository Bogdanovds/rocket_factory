package part

import (
	"github.com/bogdanovds/rocket_factory/inventory/internal/model"
	"sync"
)

type Repository struct {
	mu    sync.RWMutex
	parts map[string]*model.Part
}

func NewPartRepository() *Repository {
	return &Repository{
		parts: make(map[string]*model.Part),
	}
}
