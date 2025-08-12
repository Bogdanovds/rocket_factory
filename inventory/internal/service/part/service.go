package part

import (
	"github.com/bogdanovds/rocket_factory/inventory/internal/repository"
	"sync"
)

type Service struct {
	repo repository.PartRepository
	mu   sync.RWMutex
}

func NewPartService(repo repository.PartRepository) *Service {
	return &Service{
		repo: repo,
	}
}
