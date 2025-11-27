package mongo

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionName = "parts"
)

// Repository реализует интерфейс repository.PartRepository для MongoDB
type Repository struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

// NewRepository создаёт новый MongoDB репозиторий
func NewRepository(client *mongo.Client, dbName string) *Repository {
	db := client.Database(dbName)
	return &Repository{
		client:     client,
		database:   db,
		collection: db.Collection(collectionName),
	}
}

// Connect подключается к MongoDB
func Connect(ctx context.Context, uri string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Проверяем подключение
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("✅ Connected to MongoDB")
	return client, nil
}

// Disconnect отключается от MongoDB
func Disconnect(ctx context.Context, client *mongo.Client) error {
	if err := client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}
	return nil
}
