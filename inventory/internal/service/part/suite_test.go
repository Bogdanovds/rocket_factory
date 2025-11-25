package part

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/bogdanovds/rocket_factory/inventory/internal/repository/mocks"
)

// PartServiceTestSuite - тестовый набор для сервиса деталей
type PartServiceTestSuite struct {
	suite.Suite
	mockRepo *mocks.MockPartRepository
	service  *Service
}

// SetupTest выполняется перед каждым тестом
func (s *PartServiceTestSuite) SetupTest() {
	s.mockRepo = mocks.NewMockPartRepository()
	s.service = NewPartService(s.mockRepo)
}

// TearDownTest выполняется после каждого теста
func (s *PartServiceTestSuite) TearDownTest() {
	s.mockRepo.AssertExpectations(s.T())
}

// TestPartServiceTestSuite запускает тестовый набор
func TestPartServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PartServiceTestSuite))
}

