package mongo

import (
	"github.com/bogdanovds/rocket_factory/inventory/internal/model"
)

// ToServiceModel конвертирует документ MongoDB в модель сервисного слоя
func ToServiceModel(doc *PartDocument) *model.Part {
	if doc == nil {
		return nil
	}

	part := &model.Part{
		Uuid:          doc.UUID,
		Name:          doc.Name,
		Description:   doc.Description,
		Price:         doc.Price,
		StockQuantity: doc.StockQuantity,
		Category:      model.Category(doc.Category),
		Tags:          doc.Tags,
		Metadata:      doc.Metadata,
		CreatedAt:     &doc.CreatedAt,
		UpdatedAt:     &doc.UpdatedAt,
	}

	if doc.Dimensions != nil {
		part.Dimensions = &model.Dimensions{
			Length: doc.Dimensions.Length,
			Width:  doc.Dimensions.Width,
			Height: doc.Dimensions.Height,
			Weight: doc.Dimensions.Weight,
		}
	}

	if doc.Manufacturer != nil {
		part.Manufacturer = &model.Manufacturer{
			Name:    doc.Manufacturer.Name,
			Country: doc.Manufacturer.Country,
			Website: doc.Manufacturer.Website,
		}
	}

	return part
}

// ToDocument конвертирует модель сервисного слоя в документ MongoDB
func ToDocument(part *model.Part) *PartDocument {
	if part == nil {
		return nil
	}

	doc := &PartDocument{
		UUID:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      int32(part.Category),
		Tags:          part.Tags,
		Metadata:      part.Metadata,
	}

	if part.CreatedAt != nil {
		doc.CreatedAt = *part.CreatedAt
	}

	if part.UpdatedAt != nil {
		doc.UpdatedAt = *part.UpdatedAt
	}

	if part.Dimensions != nil {
		doc.Dimensions = &DimensionsDocument{
			Length: part.Dimensions.Length,
			Width:  part.Dimensions.Width,
			Height: part.Dimensions.Height,
			Weight: part.Dimensions.Weight,
		}
	}

	if part.Manufacturer != nil {
		doc.Manufacturer = &ManufacturerDocument{
			Name:    part.Manufacturer.Name,
			Country: part.Manufacturer.Country,
			Website: part.Manufacturer.Website,
		}
	}

	return doc
}

