package order

import (
	"github.com/bogdanovds/rocket_factory/order/internal/client"
	"github.com/bogdanovds/rocket_factory/order/internal/repository"
)

type Service struct {
	repo            repository.Repository
	inventoryClient client.InventoryClient
	paymentClient   client.PaymentClient
}

func NewService(repo repository.Repository, invClient client.InventoryClient, payClient client.PaymentClient) *Service {
	return &Service{
		repo:            repo,
		inventoryClient: invClient,
		paymentClient:   payClient,
	}
}
