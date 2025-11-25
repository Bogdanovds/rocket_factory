package order

import (
	"github.com/bogdanovds/rocket_factory/order/internal/client/grpc/inventory/v1"
	"github.com/bogdanovds/rocket_factory/order/internal/client/grpc/payment/v1"
	"github.com/bogdanovds/rocket_factory/order/internal/repository"
)

type Service struct {
	repo            repository.Repository
	inventoryClient *inventory.Client
	paymentClient   *payment.Client
}

func NewService(repo repository.Repository, invClient *inventory.Client, payClient *payment.Client) *Service {
	return &Service{
		repo:            repo,
		inventoryClient: invClient,
		paymentClient:   payClient,
	}
}
