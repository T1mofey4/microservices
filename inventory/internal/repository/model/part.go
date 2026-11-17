package model

import (
	"time"
)

type Category int32

const (
	CategoryUnspecified Category = iota
	CategoryEngine
	CategoryFuel
	CategoryPorthole
	CategoryWing
)

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

type Part struct {
	UUID          string       `bson:"uuid,omitempty"`
	Name          string       `bson:"name"`
	Description   string       `bson:"description"`
	Price         float64      `bson:"price"`
	StockQuantity int64        `bson:"stock_quantity"`
	Category      Category     `bson:"category"`
	Dimensions    Dimensions   `bson:"dimensions"`
	Manufacturer  Manufacturer `bson:"manufacturer"`
	Tags          []string     `bson:"tags"`
	CreatedAt     time.Time    `bson:"created_at"`
	UpdatedAt     *time.Time   `bson:"updated_at,omitempty"`
}
