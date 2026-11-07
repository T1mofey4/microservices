package model

import "github.com/google/uuid"

type PartInfo struct {
	PartUUID uuid.UUID
	Price    float64
}

type Category int32

const (
	CategoryUnspecified = iota
	CategoryEngine
	CategoryFuel
	CategoryPorthole
	CategoryWing
)

type PartsFilter struct {
	UUIDs                 []uuid.UUID
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}
