//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/bogdanovds/rocket_factory/order/internal/migrator"
	"github.com/bogdanovds/rocket_factory/order/internal/model"
	"github.com/bogdanovds/rocket_factory/order/internal/repository/postgres"
	tcpostgres "github.com/bogdanovds/rocket_factory/platform/pkg/testcontainers/postgres"
)

type RepositoryIntegrationTestSuite struct {
	suite.Suite
	ctx       context.Context
	container *tcpostgres.Container
	repo      *postgres.Repository
}

func (s *RepositoryIntegrationTestSuite) SetupSuite() {
	s.ctx = context.Background()

	// Создаём контейнер PostgreSQL
	container, err := tcpostgres.NewContainer(s.ctx,
		tcpostgres.WithDatabase("order_test"),
		tcpostgres.WithAuth("test", "test"),
	)
	s.Require().NoError(err)
	s.container = container

	// Применяем миграции
	m := migrator.New(container.DB())
	err = m.UpEmbed()
	s.Require().NoError(err)

	// Создаём репозиторий
	s.repo = postgres.NewRepository(container.DB())
}

func (s *RepositoryIntegrationTestSuite) TearDownSuite() {
	if s.container != nil {
		err := s.container.Terminate(s.ctx)
		s.Require().NoError(err)
	}
}

func (s *RepositoryIntegrationTestSuite) TearDownTest() {
	// Очищаем таблицу после каждого теста
	_, err := s.container.DB().ExecContext(s.ctx, "DELETE FROM orders")
	s.Require().NoError(err)
}

func (s *RepositoryIntegrationTestSuite) TestCreate_Success() {
	order := &model.Order{
		ID:            uuid.New(),
		UserID:        uuid.New(),
		PartIDs:       []uuid.UUID{uuid.New(), uuid.New()},
		TotalPrice:    150.50,
		Status:        model.OrderStatusPending,
		PaymentMethod: "CARD",
	}

	err := s.repo.Create(s.ctx, order)
	s.Require().NoError(err)

	// Проверяем, что заказ создан
	savedOrder, err := s.repo.Get(s.ctx, order.ID)
	s.Require().NoError(err)
	s.Equal(order.ID, savedOrder.ID)
	s.Equal(order.UserID, savedOrder.UserID)
	s.Equal(order.TotalPrice, savedOrder.TotalPrice)
	s.Equal(order.Status, savedOrder.Status)
	s.Equal(order.PaymentMethod, savedOrder.PaymentMethod)
	s.Equal(len(order.PartIDs), len(savedOrder.PartIDs))
}

func (s *RepositoryIntegrationTestSuite) TestCreate_WithTransactionID() {
	transactionID := uuid.New()
	order := &model.Order{
		ID:            uuid.New(),
		UserID:        uuid.New(),
		PartIDs:       []uuid.UUID{uuid.New()},
		TotalPrice:    100.00,
		Status:        model.OrderStatusPaid,
		PaymentMethod: "SBP",
		TransactionID: transactionID,
	}

	err := s.repo.Create(s.ctx, order)
	s.Require().NoError(err)

	savedOrder, err := s.repo.Get(s.ctx, order.ID)
	s.Require().NoError(err)
	s.Equal(transactionID, savedOrder.TransactionID)
}

func (s *RepositoryIntegrationTestSuite) TestGet_NotFound() {
	nonExistentID := uuid.New()

	order, err := s.repo.Get(s.ctx, nonExistentID)
	s.ErrorIs(err, model.ErrOrderNotFound)
	s.Nil(order)
}

func (s *RepositoryIntegrationTestSuite) TestGet_Success() {
	order := &model.Order{
		ID:            uuid.New(),
		UserID:        uuid.New(),
		PartIDs:       []uuid.UUID{uuid.New(), uuid.New(), uuid.New()},
		TotalPrice:    500.00,
		Status:        model.OrderStatusPending,
		PaymentMethod: "",
	}

	err := s.repo.Create(s.ctx, order)
	s.Require().NoError(err)

	savedOrder, err := s.repo.Get(s.ctx, order.ID)
	s.Require().NoError(err)

	s.Equal(order.ID, savedOrder.ID)
	s.Equal(order.UserID, savedOrder.UserID)
	s.Equal(order.TotalPrice, savedOrder.TotalPrice)
	s.Equal(order.Status, savedOrder.Status)
	s.Len(savedOrder.PartIDs, 3)

	// Проверяем, что все part IDs совпадают
	for i, partID := range order.PartIDs {
		s.Equal(partID, savedOrder.PartIDs[i])
	}
}

func (s *RepositoryIntegrationTestSuite) TestUpdate_Success() {
	order := &model.Order{
		ID:            uuid.New(),
		UserID:        uuid.New(),
		PartIDs:       []uuid.UUID{uuid.New()},
		TotalPrice:    100.00,
		Status:        model.OrderStatusPending,
		PaymentMethod: "",
	}

	err := s.repo.Create(s.ctx, order)
	s.Require().NoError(err)

	// Обновляем заказ
	order.Status = model.OrderStatusPaid
	order.PaymentMethod = "CARD"
	order.TransactionID = uuid.New()

	err = s.repo.Update(s.ctx, order)
	s.Require().NoError(err)

	// Проверяем обновление
	updatedOrder, err := s.repo.Get(s.ctx, order.ID)
	s.Require().NoError(err)
	s.Equal(model.OrderStatusPaid, updatedOrder.Status)
	s.Equal("CARD", updatedOrder.PaymentMethod)
	s.Equal(order.TransactionID, updatedOrder.TransactionID)
}

func (s *RepositoryIntegrationTestSuite) TestUpdate_NotFound() {
	order := &model.Order{
		ID:            uuid.New(),
		UserID:        uuid.New(),
		PartIDs:       []uuid.UUID{uuid.New()},
		TotalPrice:    100.00,
		Status:        model.OrderStatusPending,
		PaymentMethod: "",
	}

	err := s.repo.Update(s.ctx, order)
	s.ErrorIs(err, model.ErrOrderNotFound)
}

func (s *RepositoryIntegrationTestSuite) TestUpdate_StatusChange() {
	order := &model.Order{
		ID:            uuid.New(),
		UserID:        uuid.New(),
		PartIDs:       []uuid.UUID{uuid.New()},
		TotalPrice:    250.00,
		Status:        model.OrderStatusPending,
		PaymentMethod: "",
	}

	err := s.repo.Create(s.ctx, order)
	s.Require().NoError(err)

	// Меняем статус на Cancelled
	order.Status = model.OrderStatusCancelled
	err = s.repo.Update(s.ctx, order)
	s.Require().NoError(err)

	updatedOrder, err := s.repo.Get(s.ctx, order.ID)
	s.Require().NoError(err)
	s.Equal(model.OrderStatusCancelled, updatedOrder.Status)
}

func (s *RepositoryIntegrationTestSuite) TestCreateMultipleOrders() {
	userID := uuid.New()

	orders := []*model.Order{
		{
			ID:         uuid.New(),
			UserID:     userID,
			PartIDs:    []uuid.UUID{uuid.New()},
			TotalPrice: 100.00,
			Status:     model.OrderStatusPending,
		},
		{
			ID:         uuid.New(),
			UserID:     userID,
			PartIDs:    []uuid.UUID{uuid.New(), uuid.New()},
			TotalPrice: 200.00,
			Status:     model.OrderStatusPaid,
		},
		{
			ID:         uuid.New(),
			UserID:     uuid.New(), // Другой пользователь
			PartIDs:    []uuid.UUID{uuid.New()},
			TotalPrice: 50.00,
			Status:     model.OrderStatusCancelled,
		},
	}

	for _, order := range orders {
		err := s.repo.Create(s.ctx, order)
		s.Require().NoError(err)
	}

	// Проверяем каждый заказ
	for _, order := range orders {
		savedOrder, err := s.repo.Get(s.ctx, order.ID)
		s.Require().NoError(err)
		s.Equal(order.ID, savedOrder.ID)
		s.Equal(order.Status, savedOrder.Status)
	}
}

func TestRepositoryIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryIntegrationTestSuite))
}

