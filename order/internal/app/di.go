package app

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	v1 "github.com/bogdanovds/rocket_factory/order/internal/api/order/v1"
	"github.com/bogdanovds/rocket_factory/order/internal/client"
	inventoryClient "github.com/bogdanovds/rocket_factory/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/bogdanovds/rocket_factory/order/internal/client/grpc/payment/v1"
	"github.com/bogdanovds/rocket_factory/order/internal/config"
	"github.com/bogdanovds/rocket_factory/order/internal/migrator"
	"github.com/bogdanovds/rocket_factory/order/internal/repository"
	"github.com/bogdanovds/rocket_factory/order/internal/repository/postgres"
	"github.com/bogdanovds/rocket_factory/order/internal/service"
	orderService "github.com/bogdanovds/rocket_factory/order/internal/service/order"
	"github.com/bogdanovds/rocket_factory/platform/pkg/closer"
	orderV1 "github.com/bogdanovds/rocket_factory/shared/pkg/openapi/order/v1"
)

type diContainer struct {
	orderV1Handler orderV1.Handler

	orderService service.Service

	orderRepository repository.Repository

	inventoryClient client.InventoryClient
	paymentClient   client.PaymentClient

	inventoryConn *grpc.ClientConn
	paymentConn   *grpc.ClientConn

	db *sql.DB
}

func newDIContainer() *diContainer {
	return &diContainer{}
}

// OrderV1Handler возвращает HTTP handler
func (d *diContainer) OrderV1Handler(ctx context.Context) orderV1.Handler {
	if d.orderV1Handler == nil {
		d.orderV1Handler = v1.NewHandler(d.OrderService(ctx))
	}

	return d.orderV1Handler
}

// OrderService возвращает сервис заказов
func (d *diContainer) OrderService(ctx context.Context) service.Service {
	if d.orderService == nil {
		d.orderService = orderService.NewService(
			d.OrderRepository(ctx),
			d.InventoryClient(ctx),
			d.PaymentClient(ctx),
		)
	}

	return d.orderService
}

// OrderRepository возвращает репозиторий заказов
func (d *diContainer) OrderRepository(ctx context.Context) repository.Repository {
	if d.orderRepository == nil {
		d.orderRepository = postgres.NewRepository(d.DB(ctx))
	}

	return d.orderRepository
}

// InventoryClient возвращает клиент Inventory
func (d *diContainer) InventoryClient(ctx context.Context) client.InventoryClient {
	if d.inventoryClient == nil {
		d.inventoryClient = inventoryClient.New(d.InventoryGRPCConn(ctx))
	}

	return d.inventoryClient
}

// PaymentClient возвращает клиент Payment
func (d *diContainer) PaymentClient(ctx context.Context) client.PaymentClient {
	if d.paymentClient == nil {
		d.paymentClient = paymentClient.New(d.PaymentGRPCConn(ctx))
	}

	return d.paymentClient
}

// InventoryGRPCConn возвращает gRPC соединение с Inventory
func (d *diContainer) InventoryGRPCConn(ctx context.Context) *grpc.ClientConn {
	if d.inventoryConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().InventoryClient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to inventory gRPC: %v", err))
		}

		closer.AddNamed("Inventory gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.inventoryConn = conn
	}

	return d.inventoryConn
}

// PaymentGRPCConn возвращает gRPC соединение с Payment
func (d *diContainer) PaymentGRPCConn(ctx context.Context) *grpc.ClientConn {
	if d.paymentConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().PaymentClient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to payment gRPC: %v", err))
		}

		closer.AddNamed("Payment gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.paymentConn = conn
	}

	return d.paymentConn
}

// DB возвращает соединение с PostgreSQL
func (d *diContainer) DB(ctx context.Context) *sql.DB {
	if d.db == nil {
		db, err := sql.Open("postgres", config.AppConfig().Postgres.DSN())
		if err != nil {
			panic(fmt.Sprintf("failed to open database: %v", err))
		}

		if err = db.PingContext(ctx); err != nil {
			panic(fmt.Sprintf("failed to ping database: %v", err))
		}

		// Применяем миграции
		m := migrator.New(db)
		if err = m.UpEmbed(); err != nil {
			panic(fmt.Sprintf("failed to apply migrations: %v", err))
		}

		closer.AddNamed("PostgreSQL connection", func(ctx context.Context) error {
			return db.Close()
		})

		d.db = db
	}

	return d.db
}
