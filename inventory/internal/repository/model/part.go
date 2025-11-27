package model

import "time"

// Part - модель детали для слоя repository
type Part struct {
	Uuid          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      int32
	Length        float64
	Width         float64
	Height        float64
	Weight        float64
	Manufacturer  string
	Country       string
	Website       string
	Tags          []string
	Metadata      map[string]interface{}
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}
