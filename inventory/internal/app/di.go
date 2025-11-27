package app

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"

	api "github.com/bogdanovds/rocket_factory/inventory/internal/api/inventory/v1"
	"github.com/bogdanovds/rocket_factory/inventory/internal/config"
	"github.com/bogdanovds/rocket_factory/inventory/internal/repository"
	mongoRepo "github.com/bogdanovds/rocket_factory/inventory/internal/repository/mongo"
	"github.com/bogdanovds/rocket_factory/inventory/internal/service"
	partService "github.com/bogdanovds/rocket_factory/inventory/internal/service/part"
	"github.com/bogdanovds/rocket_factory/platform/pkg/closer"
	"github.com/bogdanovds/rocket_factory/platform/pkg/logger"
	inventoryV1 "github.com/bogdanovds/rocket_factory/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	inventoryV1API inventoryV1.InventoryServiceServer

	partService service.PartService

	partRepository repository.PartRepository

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database
}

func newDIContainer() *diContainer {
	return &diContainer{}
}

// InventoryV1API возвращает gRPC API сервер
func (d *diContainer) InventoryV1API(ctx context.Context) inventoryV1.InventoryServiceServer {
	if d.inventoryV1API == nil {
		d.inventoryV1API = api.NewInventoryAPI(d.PartService(ctx))
	}

	return d.inventoryV1API
}

// PartService возвращает сервис деталей
func (d *diContainer) PartService(ctx context.Context) service.PartService {
	if d.partService == nil {
		d.partService = partService.NewPartService(d.PartRepository(ctx))
	}

	return d.partService
}

// PartRepository возвращает репозиторий деталей
func (d *diContainer) PartRepository(ctx context.Context) repository.PartRepository {
	if d.partRepository == nil {
		repo := mongoRepo.NewRepository(d.MongoDBClient(ctx), config.AppConfig().Mongo.DatabaseName())

		// Заполняем начальные данные
		if err := repo.SeedParts(ctx); err != nil {
			// Логируем, но не падаем - данные могут уже существовать
			logger.Warn(ctx, "Failed to seed parts", zap.Error(err))
		}

		d.partRepository = repo
	}

	return d.partRepository
}

// MongoDBClient возвращает клиент MongoDB
func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %v\n", err))
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		d.mongoDBClient = client
	}

	return d.mongoDBClient
}

// MongoDBHandle возвращает handle базы данных MongoDB
func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
	}

	return d.mongoDBHandle
}
