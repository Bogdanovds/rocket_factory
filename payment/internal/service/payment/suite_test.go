package payment

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// PaymentServiceTestSuite - тестовый набор для сервиса оплаты
type PaymentServiceTestSuite struct {
	suite.Suite
	service *Service
}

// SetupTest выполняется перед каждым тестом
func (s *PaymentServiceTestSuite) SetupTest() {
	s.service = NewPaymentService()
}

// TestPaymentServiceTestSuite запускает тестовый набор
func TestPaymentServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentServiceTestSuite))
}
