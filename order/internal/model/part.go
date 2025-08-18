package model

import "github.com/google/uuid"

type Part struct {
	ID       uuid.UUID
	Name     string
	Price    float64
	Category string
}
