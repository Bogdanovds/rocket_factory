package order

import (
	"testing"

	"github.com/stretchr/testify/suite"

	clientMocks "github.com/bogdanovds/rocket_factory/order/internal/client/grpc/mocks"
	repoMocks "github.com/bogdanovds/rocket_factory/order/internal/repository/mocks"
)

// OrderServiceTestSuite - тестовый набор для сервиса заказов
type OrderServiceTestSuite struct {
	suite.Suite
	mockRepo            *repoMocks.MockOrderRepository
	mockInventoryClient *clientMocks.MockInventoryClient
	mockPaymentClient   *clientMocks.MockPaymentClient
	service             *Service
}

// SetupTest выполняется перед каждым тестом
func (s *OrderServiceTestSuite) SetupTest() {
	s.mockRepo = repoMocks.NewMockOrderRepository()
	s.mockInventoryClient = clientMocks.NewMockInventoryClient()
	s.mockPaymentClient = clientMocks.NewMockPaymentClient()
	s.service = NewService(s.mockRepo, s.mockInventoryClient, s.mockPaymentClient)
}

// TearDownTest выполняется после каждого теста
func (s *OrderServiceTestSuite) TearDownTest() {
	s.mockRepo.AssertExpectations(s.T())
	s.mockInventoryClient.AssertExpectations(s.T())
	s.mockPaymentClient.AssertExpectations(s.T())
}

// TestOrderServiceTestSuite запускает тестовый набор
func TestOrderServiceTestSuite(t *testing.T) {
	suite.Run(t, new(OrderServiceTestSuite))
}

