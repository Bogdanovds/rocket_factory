package converter

import (
	"github.com/bogdanovds/rocket_factory/inventory/internal/model"
	repoModel "github.com/bogdanovds/rocket_factory/inventory/internal/repository/model"
)

// ToServiceModel конвертирует модель repository в модель сервисного слоя
func ToServiceModel(part *repoModel.Part) *model.Part {
	if part == nil {
		return nil
	}

	return &model.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      model.Category(part.Category),
		Dimensions: &model.Dimensions{
			Length: part.Length,
			Width:  part.Width,
			Height: part.Height,
			Weight: part.Weight,
		},
		Manufacturer: &model.Manufacturer{
			Name:    part.Manufacturer,
			Country: part.Country,
			Website: part.Website,
		},
		Tags:      part.Tags,
		Metadata:  part.Metadata,
		CreatedAt: part.CreatedAt,
		UpdatedAt: part.UpdatedAt,
	}
}

// ToRepoModel конвертирует модель сервисного слоя в модель repository
func ToRepoModel(part *model.Part) *repoModel.Part {
	if part == nil {
		return nil
	}

	repoPart := &repoModel.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      int32(part.Category),
		Tags:          part.Tags,
		Metadata:      part.Metadata,
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}

	if part.Dimensions != nil {
		repoPart.Length = part.Dimensions.Length
		repoPart.Width = part.Dimensions.Width
		repoPart.Height = part.Dimensions.Height
		repoPart.Weight = part.Dimensions.Weight
	}

	if part.Manufacturer != nil {
		repoPart.Manufacturer = part.Manufacturer.Name
		repoPart.Country = part.Manufacturer.Country
		repoPart.Website = part.Manufacturer.Website
	}

	return repoPart
}

// ToServiceModelList конвертирует список моделей repository в список моделей сервисного слоя
func ToServiceModelList(parts []*repoModel.Part) []*model.Part {
	result := make([]*model.Part, 0, len(parts))
	for _, p := range parts {
		result = append(result, ToServiceModel(p))
	}
	return result
}
