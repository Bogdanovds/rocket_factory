package mongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/bogdanovds/rocket_factory/inventory/internal/model"
)

// SeedParts заполняет MongoDB начальными данными
func (r *Repository) SeedParts(ctx context.Context) error {
	// Проверяем, есть ли уже данные
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}

	if count > 0 {
		log.Printf("ℹ️ Parts collection already has %d documents, skipping seed", count)
		return nil
	}

	now := time.Now()

	parts := []*model.Part{
		{
			Uuid:          "6ba7b810-9dad-11d1-80b4-00c04fd430c9",
			Name:          "Main Engine",
			Description:   "Primary propulsion system",
			Price:         2500000.99,
			StockQuantity: 5,
			Category:      model.CategoryEngine,
			Dimensions:    &model.Dimensions{Length: 450, Width: 200, Height: 300, Weight: 8500},
			Manufacturer:  &model.Manufacturer{Name: "SpaceTech", Country: "USA", Website: "spacetech.com"},
			Tags:          []string{"propulsion", "primary", "engine"},
			CreatedAt:     &now,
			UpdatedAt:     &now,
		},
		{
			Uuid:          "6ba7b810-9dad-11d1-80b4-00c04fd430ca",
			Name:          "Fuel Tank",
			Description:   "Liquid hydrogen storage",
			Price:         1200000.50,
			StockQuantity: 8,
			Category:      model.CategoryFuel,
			Dimensions:    &model.Dimensions{Length: 600, Width: 300, Height: 300, Weight: 2000},
			Manufacturer:  &model.Manufacturer{Name: "FuelSystems", Country: "Germany", Website: "fuelsystems.de"},
			Tags:          []string{"storage", "fuel", "hydrogen"},
			CreatedAt:     &now,
			UpdatedAt:     &now,
		},
		{
			Uuid:          "6ba7b810-9dad-11d1-80b4-00c04fd430cb",
			Name:          "Navigation Computer",
			Description:   "Advanced flight navigation system",
			Price:         850000.00,
			StockQuantity: 12,
			Category:      model.CategoryEngine,
			Dimensions:    &model.Dimensions{Length: 50, Width: 40, Height: 30, Weight: 25},
			Manufacturer:  &model.Manufacturer{Name: "NavTech", Country: "Japan", Website: "navtech.jp"},
			Tags:          []string{"navigation", "computer", "avionics"},
			CreatedAt:     &now,
			UpdatedAt:     &now,
		},
		{
			Uuid:          "6ba7b810-9dad-11d1-80b4-00c04fd430cc",
			Name:          "Heat Shield",
			Description:   "Thermal protection system for reentry",
			Price:         1800000.00,
			StockQuantity: 3,
			Category:      model.CategoryWing,
			Dimensions:    &model.Dimensions{Length: 800, Width: 600, Height: 100, Weight: 1500},
			Manufacturer:  &model.Manufacturer{Name: "ThermalCorp", Country: "USA", Website: "thermalcorp.com"},
			Tags:          []string{"thermal", "protection", "reentry"},
			CreatedAt:     &now,
			UpdatedAt:     &now,
		},
	}

	docs := make([]interface{}, len(parts))
	for i, part := range parts {
		docs[i] = ToDocument(part)
	}

	_, err = r.collection.InsertMany(ctx, docs)
	if err != nil {
		return err
	}

	log.Printf("✅ Seeded %d parts into MongoDB", len(parts))
	return nil
}
