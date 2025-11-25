package part

import (
	"sync"

	"github.com/bogdanovds/rocket_factory/inventory/internal/repository"
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
