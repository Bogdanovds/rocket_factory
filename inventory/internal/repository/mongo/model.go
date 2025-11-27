package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PartDocument - структура документа в MongoDB
type PartDocument struct {
	ID            primitive.ObjectID     `bson:"_id,omitempty"`
	UUID          string                 `bson:"uuid"`
	Name          string                 `bson:"name"`
	Description   string                 `bson:"description"`
	Price         float64                `bson:"price"`
	StockQuantity int64                  `bson:"stock_quantity"`
	Category      int32                  `bson:"category"`
	Dimensions    *DimensionsDocument    `bson:"dimensions,omitempty"`
	Manufacturer  *ManufacturerDocument  `bson:"manufacturer,omitempty"`
	Tags          []string               `bson:"tags"`
	Metadata      map[string]interface{} `bson:"metadata,omitempty"`
	CreatedAt     time.Time              `bson:"created_at"`
	UpdatedAt     time.Time              `bson:"updated_at"`
}

// DimensionsDocument - структура размеров
type DimensionsDocument struct {
	Length float64 `bson:"length"`
	Width  float64 `bson:"width"`
	Height float64 `bson:"height"`
	Weight float64 `bson:"weight"`
}

// ManufacturerDocument - структура производителя
type ManufacturerDocument struct {
	Name    string `bson:"name"`
	Country string `bson:"country"`
	Website string `bson:"website"`
}
