package converter

import (
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/bogdanovds/rocket_factory/inventory/internal/model"
	inventoryV1 "github.com/bogdanovds/rocket_factory/shared/pkg/proto/inventory/v1"
)

func ToProtoPart(p *model.Part) *inventoryV1.Part {
	metadata := make(map[string]*inventoryV1.Value)
	for k, v := range p.Metadata {
		metadata[k] = toProtoValue(v)
	}

	var createdAt, updatedAt *timestamppb.Timestamp
	if p.CreatedAt != nil {
		createdAt = timestamppb.New(*p.CreatedAt)
	}
	if p.UpdatedAt != nil {
		updatedAt = timestamppb.New(*p.UpdatedAt)
	}

	return &inventoryV1.Part{
		Uuid:          p.Uuid,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		Category:      inventoryV1.Category(p.Category),
		Dimensions: &inventoryV1.Dimensions{
			Length: p.Dimensions.Length,
			Width:  p.Dimensions.Width,
			Height: p.Dimensions.Height,
			Weight: p.Dimensions.Weight,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    p.Manufacturer.Name,
			Country: p.Manufacturer.Country,
			Website: p.Manufacturer.Website,
		},
		Tags:      p.Tags,
		Metadata:  metadata,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func ToModelPart(p *inventoryV1.Part) *model.Part {
	metadata := make(map[string]interface{})
	for k, v := range p.GetMetadata() {
		metadata[k] = toModelValue(v)
	}

	return &model.Part{
		Uuid:          p.GetUuid(),
		Name:          p.GetName(),
		Description:   p.GetDescription(),
		Price:         p.GetPrice(),
		StockQuantity: p.GetStockQuantity(),
		Category:      model.Category(p.GetCategory()),
		Dimensions: &model.Dimensions{
			Length: p.GetDimensions().GetLength(),
			Width:  p.GetDimensions().GetWidth(),
			Height: p.GetDimensions().GetHeight(),
			Weight: p.GetDimensions().GetWeight(),
		},
		Manufacturer: &model.Manufacturer{
			Name:    p.GetManufacturer().GetName(),
			Country: p.GetManufacturer().GetCountry(),
			Website: p.GetManufacturer().GetWebsite(),
		},
		Tags:      p.GetTags(),
		Metadata:  metadata,
		CreatedAt: lo.ToPtr(p.GetCreatedAt().AsTime()),
		UpdatedAt: lo.ToPtr(p.GetUpdatedAt().AsTime()),
	}
}

func toProtoValue(v interface{}) *inventoryV1.Value {
	switch val := v.(type) {
	case string:
		return &inventoryV1.Value{Value: &inventoryV1.Value_StringValue{StringValue: val}}
	case int64:
		return &inventoryV1.Value{Value: &inventoryV1.Value_Int64Value{Int64Value: val}}
	case float64:
		return &inventoryV1.Value{Value: &inventoryV1.Value_DoubleValue{DoubleValue: val}}
	case bool:
		return &inventoryV1.Value{Value: &inventoryV1.Value_BoolValue{BoolValue: val}}
	default:
		return nil
	}
}

func toModelValue(v *inventoryV1.Value) interface{} {
	switch val := v.Value.(type) {
	case *inventoryV1.Value_StringValue:
		return val.StringValue
	case *inventoryV1.Value_Int64Value:
		return val.Int64Value
	case *inventoryV1.Value_DoubleValue:
		return val.DoubleValue
	case *inventoryV1.Value_BoolValue:
		return val.BoolValue
	default:
		return nil
	}
}

func ToProtoPartsFilter(f *inventoryV1.PartsFilter) *model.PartsFilter {
	if f == nil {
		return nil
	}

	categories := make([]model.Category, 0, len(f.GetCategories()))
	for _, c := range f.GetCategories() {
		categories = append(categories, model.Category(c))
	}

	return &model.PartsFilter{
		Uuids:                 f.GetUuids(),
		Names:                 f.GetNames(),
		Categories:            categories,
		ManufacturerCountries: f.GetManufacturerCountries(),
		Tags:                  f.GetTags(),
	}
}
