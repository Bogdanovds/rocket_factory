package part

import (
	"github.com/bogdanovds/rocket_factory/inventory/internal/model"
	"github.com/samber/lo"
	"time"

	"github.com/bogdanovds/rocket_factory/inventory/internal/repository"
)

func SeedParts(repo repository.PartRepository) {
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
			CreatedAt:     lo.ToPtr(now),
			UpdatedAt:     lo.ToPtr(now),
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
			CreatedAt:     lo.ToPtr(now),
			UpdatedAt:     lo.ToPtr(now),
		},
	}

	if repoImpl, ok := repo.(*Repository); ok {
		for _, part := range parts {
			repoImpl.parts[part.Uuid] = part
		}
	}
}
