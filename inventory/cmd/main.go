package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	api "github.com/bogdanovds/rocket_factory/inventory/internal/api/inventory/v1"
	mongoRepo "github.com/bogdanovds/rocket_factory/inventory/internal/repository/mongo"
	service "github.com/bogdanovds/rocket_factory/inventory/internal/service/part"
)

const grpcPort = 50051

func main() {
	ctx := context.Background()

	// Подключение к MongoDB
	mongoURI := getEnv("MONGO_URI", "mongodb://inventory-service-user:inventory-service-password@localhost:27017")
	mongoDBName := getEnv("MONGO_DB", "inventory-service")

	mongoClient, err := mongoRepo.Connect(ctx, mongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := mongoRepo.Disconnect(ctx, mongoClient); err != nil {
			log.Printf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	partRepo := mongoRepo.NewRepository(mongoClient, mongoDBName)

	// Заполняем начальные данные
	if err := partRepo.SeedParts(ctx); err != nil {
		log.Printf("Warning: Failed to seed parts: %v", err)
	}

	partService := service.NewPartService(partRepo)
	inventoryAPI := api.NewInventoryAPI(partService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	s := grpc.NewServer()
	api.RegisterInventoryServiceServer(s, inventoryAPI)
	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
