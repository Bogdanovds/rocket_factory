package main

import (
	"fmt"
	repo "github.com/bogdanovds/rocket_factory/inventory/internal/repository/part"
	service "github.com/bogdanovds/rocket_factory/inventory/internal/service/part"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	api "github.com/bogdanovds/rocket_factory/inventory/internal/api/inventory/v1"
)

const grpcPort = 50051

func main() {
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

	partRepo := repo.NewPartRepository()
	partService := service.NewPartService(partRepo)
	inventoryAPI := api.NewInventoryAPI(partService)

	s := grpc.NewServer()
	api.RegisterInventoryServiceServer(s, inventoryAPI)
	reflection.Register(s)

	repo.SeedParts(partRepo)

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
