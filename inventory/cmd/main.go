package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	inventoryV1 "github.com/bogdanovds/rocket_factory/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

type InventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer
	mu    sync.RWMutex
	parts map[string]*inventoryV1.Part
}

func NewInventoryService() *InventoryService {
	return &InventoryService{
		parts: make(map[string]*inventoryV1.Part),
	}
}

func (s *InventoryService) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	part, exists := s.parts[req.Uuid]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.Uuid)
	}

	return &inventoryV1.GetPartResponse{Part: part}, nil
}

func (s *InventoryService) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var parts []*inventoryV1.Part

	// Если фильтр пустой - возвращаем все детали
	if req.Filter == nil || isEmptyFilter(req.Filter) {
		parts = make([]*inventoryV1.Part, 0, len(s.parts))
		for _, part := range s.parts {
			parts = append(parts, part)
		}
		return &inventoryV1.ListPartsResponse{Parts: parts}, nil
	}

	// Фильтрация деталей
	for _, part := range s.parts {
		if matchesFilter(part, req.Filter) {
			parts = append(parts, part)
		}
	}

	return &inventoryV1.ListPartsResponse{Parts: parts}, nil
}

func isEmptyFilter(filter *inventoryV1.PartsFilter) bool {
	return len(filter.Uuids) == 0 &&
		len(filter.Names) == 0 &&
		len(filter.Categories) == 0 &&
		len(filter.ManufacturerCountries) == 0 &&
		len(filter.Tags) == 0
}

func matchesFilter(part *inventoryV1.Part, filter *inventoryV1.PartsFilter) bool {
	// Фильтрация по UUID
	if len(filter.Uuids) > 0 && !contains(filter.Uuids, part.Uuid) {
		return false
	}

	// Фильтрация по имени
	if len(filter.Names) > 0 && !contains(filter.Names, part.Name) {
		return false
	}

	// Фильтрация по категории
	if len(filter.Categories) > 0 && !containsCategory(filter.Categories, part.Category) {
		return false
	}

	// Фильтрация по стране производителя
	if len(filter.ManufacturerCountries) > 0 && !contains(filter.ManufacturerCountries, part.Manufacturer.Country) {
		return false
	}

	// Фильтрация по тегам
	if len(filter.Tags) > 0 && !hasAnyTag(part.Tags, filter.Tags) {
		return false
	}

	return true
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func containsCategory(categories []inventoryV1.Category, category inventoryV1.Category) bool {
	for _, c := range categories {
		if c == category {
			return true
		}
	}
	return false
}

func hasAnyTag(partTags, filterTags []string) bool {
	for _, ft := range filterTags {
		for _, pt := range partTags {
			if pt == ft {
				return true
			}
		}
	}
	return false
}

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

	s := grpc.NewServer()
	inventoryService := NewInventoryService()
	inventoryV1.RegisterInventoryServiceServer(s, inventoryService)

	reflection.Register(s)

	// Добавление тестовых данных
	seedInventoryService(inventoryService)

	go func() {
		log.Printf("🚀 gRPC s listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC s...")
	s.GracefulStop()
	log.Println("✅ s stopped")
}

func seedInventoryService(service *InventoryService) {
	now := timestamppb.Now()

	parts := []*inventoryV1.Part{
		{
			Uuid:          "6ba7b810-9dad-11d1-80b4-00c04fd430c9",
			Name:          "Main Engine",
			Description:   "Primary propulsion system",
			Price:         2500000.99,
			StockQuantity: 5,
			Category:      inventoryV1.Category_CATEGORY_ENGINE,
			Dimensions:    &inventoryV1.Dimensions{Length: 450, Width: 200, Height: 300, Weight: 8500},
			Manufacturer:  &inventoryV1.Manufacturer{Name: "SpaceTech", Country: "USA", Website: "spacetech.com"},
			Tags:          []string{"propulsion", "primary", "engine"},
			CreatedAt:     now,
			UpdatedAt:     now,
		},
		{
			Uuid:          "6ba7b810-9dad-11d1-80b4-00c04fd430ca",
			Name:          "Fuel Tank",
			Description:   "Liquid hydrogen storage",
			Price:         1200000.50,
			StockQuantity: 8,
			Category:      inventoryV1.Category_CATEGORY_FUEL,
			Dimensions:    &inventoryV1.Dimensions{Length: 600, Width: 300, Height: 300, Weight: 2000},
			Manufacturer:  &inventoryV1.Manufacturer{Name: "FuelSystems", Country: "Germany", Website: "fuelsystems.de"},
			Tags:          []string{"storage", "fuel", "hydrogen"},
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}

	service.mu.Lock()
	defer service.mu.Unlock()

	for _, part := range parts {
		service.parts[part.Uuid] = part
	}
}
